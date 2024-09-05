/* drwbases: win.rexx */

say  '[win] direct entry'
exit 'FAIL'

/* ------------------------------------------------------------------ */

$file_exists: procedure
  if stream( arg( 1 ), 'C', 'QUERY EXISTS' ) ^= '' then return 1
  return 0

/* ------------------------------------------------------------------ */

$instdir_home: procedure
  filespec = arg( 1 )
  parse var filespec . '/' fname
  call RegGetVal( 'HLM',,
                  'SOFTWARE\IDAVLab\Enterprise Suite\Dr.Web (R) Enterprise Agent\Settings',,
                  'Home', st )
  return st.1'\'fname

/* ------------------------------------------------------------------ */

$reset_component_if_failed : procedure
  c = arg(1)
  if c = '' then return
  
  say '[reset failed component] process component 'c
  /* open key */
  call ApiSigErr 'ARG'
  hk = RegOpen( 'HLM', 'SOFTWARE\IDAVLab\Enterprise Suite\Dr.Web (R) Network Installer\Products\'c ) /* ' */
  if hk ^= 'FAILED' then do
    /* read state */
    err = RegGetVal(hk,,,v)

    if err = 'OK' then do
      if find(v.1, 'FAIL') ^= 0 then do
        /* remove records about every revision except previous */
        prevrev = word(v.1, 1)
        say '[reset failed component] fall back to previous working revision 'prevrev
        call RegEnumVal hk, 'kn', 'kt', 'kv'
        do I = 1 to kn.0
          if kn.I = '' | kn.I = prevrev then iterate
          say '[reset failed component] delete 'kn.I' revision info'
          call RegSetVal hk,,kn.I
        end
        call RegSetVal( hk,,, 'REG_SZ', prevrev )
      end      
    end
  
    call RegClose(hk)
  end
  call ApiSigErr 'API'
  return 
  
/* Check whether we must perform some operations on set of files. 
 * Return 1 if we must. Usually we pass here only one file
 * First it fill set of compounds variables with stems s_equ, s_add, s_del 
 * and s_rep using UpgDiffState defined in script-api.cxx
 * Every compaund variable with tail '0' contains number of files that 
 * corresponding operation must handle. For example, after executing 
 * UpgDiffState number of files which must be replaced is stored in s_rep.0
 * 
 * First param - string with concatenated operations names that we want to 
 * perform on set of files. For example 'DelAddRep'. Four operations are 
 * allowed 'Rep', 'Add', 'Del', 'Equ'.
 *
 * Second and other params contains files names. We check whether some 
 * operations can be performed for them */
$fileset: procedure
  call UpgDiffState 's_equ','s_del','s_add','s_rep'
  oper = arg( 1 )
  do I = 2 to arg()
    if pos( 'Rep', oper ) ^= 0 then do J = 1 to s_rep.0
      if translate( subword( s_rep.J, 3 ) ) = translate( arg( I ) ) then return 1
    end
    if pos( 'Add', oper ) ^= 0 then do J = 1 to s_add.0
      if translate( subword( s_add.J, 3 ) ) = translate( arg( I ) ) then return 1
    end
    if pos( 'Del', oper ) ^= 0 then do J = 1 to s_del.0
      if translate( subword( s_del.J, 3 ) ) = translate( arg( I ) ) then return 1
    end
    if pos( 'Equ', oper ) ^= 0 then do J = 1 to s_equ.0
      if translate( subword( s_equ.J, 3 ) ) = translate( arg( I ) ) then return 1
    end
  end
  return 0

/* First run special fixup script if it is added to 10-drwbases or was changed */  
$run_fixup: procedure
  fixup = UpgBase() || UpgProduct() || '\common\FIXUP.REXX'
  r1 = $fileset('AddRep', 'common/FIXUP.REXX')
  r2 = $file_exists(fixup)
  if r1 & r2 then do
    if stream(fixup, 'C', 'OPEN READ') = 0 then do
      say 'unable to open <'fixup'> for reading'
      return
    end
    
    /* read script */
    fixup_scr = charin(fixup, 1, stream(fixup, 'C', 'QUERY SIZE'))
    call stream(fixup, 'C', 'CLOSE')
    
    /* run it */
    say 'run fixup script <'fixup'>'
    
    signal off syntax
    signal off error
    signal off failure    
    interpret fixup_scr
    signal on syntax  name log_sig
    signal on error   name log_sig
    signal on failure name log_sig    
  end
  else say 'there is no new fixup script... doesn`t need any hotfixes'
  return

$do_upgrade:
  ret = 'OK'

  call $run_fixup
  
  if $checkpoint( 1, ret ) then ret = $files( $before_files( ret ) )
  if $checkpoint( 2, ret ) then return $after_files( ret )
  say 'INVALID CHECKPOINT STATE'
  exit 'FAIL'

$checkpoint:
  stage = arg( 1 )
  ret   = arg( 2 )
  if ret = 'DELAY' then do
    say '[checkpoint] delay stage 'stage
    exit ret || '.' || stage
  end
  if UpgStage() > stage then do
    say '[checkpoint] skip stage 'stage' (frozen at 'UpgStage()')'
    return 0
  end
  say '[checkpoint] begin stage 'stage
  return 1

/* ------------------------------------------------------------------ */

$files:
  ret = arg( 1 )

  cachedir = UpgBase() || UpgProduct() || '/'
  flag = 0

  /* add/rep/del && verify */
  call UpgDiffState ,'s_del',,,'s_new'
  do I = 1 to s_del.0
    parse var s_del.I . . fspec
    dst = $instdir( fspec )
    say '[files] remove <'dst'>'
    call InstFile dst
  end

  /*If this is force update from 4.33 or 4.44 to 5.0 then bases will remain undeleted*/
  signal off error

  if api.file.delete_files_by_pattern = 'yes' then do
    mask433 = $instdir_home()'d??433*'
    mask444 = $instdir_home()'d??444*'
    mask500 = $instdir_home()'d??500*'
    call DeleteFilesByPattern(mask433)
    call DeleteFilesByPattern(mask444)
    if $is_spiderg3_possible() then do
      call DeleteFilesByPattern(mask500)
    end
  end
  else do
    mask433 = '"'$instdir_home()'d??433*"'
    mask444 = '"'$instdir_home()'d??444*"'
    mask500 = '"'$instdir_home()'d??500*"'
    cmd_line433 = 'if exist 'mask433' del 'mask433  
    cmd_line444 = 'if exist 'mask444' del 'mask444
    cmd_line500 = 'if exist 'mask500' del 'mask500
    say 'Run <'cmd_line433'> to clean up old 433 bases'
    address system cmd_line433
    say 'Run <'cmd_line444'> to clean up old 444 bases'
    address system cmd_line444
    if $is_spiderg3_possible() then do
      say 'Run <'cmd_line500'> to clean up old 500 bases'
      address system cmd_line500
    end
  end  
  
  signal on error name log_sig

  did_something = 0

  do I = 1 to s_new.0
    parse var s_new.I len md5 fspec
    if translate( fspec ) = translate( 'common/DRWTODAY.VDB' ) then do
      flag = I
    end
    else do
      src = cachedir || fspec
      dst = $instdir( fspec )
      say '[files] copy&check <'src'> to <'dst'> length 'len' digest 'md5
      call InstFile dst, len, md5, src
      did_something = 1
    end
  end

  /* should be last changed file
     forgive me: that's not correct if FlagFile arg changed in .ini :( */
  if flag ^= 0 then do
    parse var s_new.flag len md5 fspec
    src = cachedir || fspec
    dst = $instdir( fspec )
    say '[files] flag copy&check <'src'> to <'dst'> length 'len' digest 'md5
    call InstFile dst, len, md5, src
    did_something = 1
  end

  if InstDelayUntilReboot == 1 then ret = 'DELAY'
  if api.dwse.restart = 'yes' & did_something ^= 0 then call DwSERestart

  return ret;

/* ------------------------------------------------------------------ */

$before_files_win:
  return arg( 1 )

$gethomedir: procedure
  call RegGetVal( 'HLM',,
                  'SOFTWARE\IDAVLab\Enterprise Suite\Dr.Web (R) Enterprise Agent\Settings',,
                  'Home', st )
  return st.1 

$after_files_win:
  /* making a cleanup */
  home = $gethomedir()
  say 'Cleanup in 'home'\cache\10-drwbases'

  if api.file.clear_dir = 'yes' then do
    call ClearDir home'\cache\10-drwbases\win'
    call ClearDir home'\cache\10-drwbases\win-nt'
    call ClearDir home'\cache\10-drwbases\win-nt-4'
    call ClearDir home'\cache\10-drwbases\win-nt-32'
    call ClearDir home'\cache\10-drwbases\win-nt-64'
    call ClearDir home'\cache\10-drwbases\win-9x'
    call ClearDir home'\cache\10-drwbases\common'
  end
  else do
    call ApiSigErr 'ARG'
    call RemoveDir home'\cache\10-drwbases\win', 'Recurse'
    call RemoveDir home'\cache\10-drwbases\win-nt', 'Recurse'
    call RemoveDir home'\cache\10-drwbases\win-nt-4', 'Recurse'
    call RemoveDir home'\cache\10-drwbases\win-nt-32', 'Recurse'
    call RemoveDir home'\cache\10-drwbases\win-nt-64', 'Recurse'
    call RemoveDir home'\cache\10-drwbases\win-9x', 'Recurse'
    call RemoveDir home'\cache\10-drwbases\common', 'Recurse'
    call ApiSigErr 'API'
  end

  return arg( 1 )

