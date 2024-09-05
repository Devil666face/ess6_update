/* drwbases: common.rexx */

say  '[common] direct entry'
exit 'FAIL'

/* ------------------------------------------------------------------ */

log_sig:
  say '  #---------CATCH---------'
  say '  #line:  <'sigl'><'sourceline(sigl)'>'
  say '  #error: <'rc'><'errortext(rc)'>'
  say '  #condition:'
  say '  #  I<'condition( 'I' )'>'
  say '  #  C<'condition( 'C' )'>'
  say '  #  D<'condition( 'D' )'>'
  say '  #  S<'condition( 'S' )'>'
  if rc == 40 then  say '  #apierr: <'apierr'>'

  if symbol( 'sig_catch' ) ^== 'LIT' & sig_catch ^= '' then do
    signal on syntax  name log_syn
    signal on error   name log_err
    signal on failure name log_fail
    sav_catch = sig_catch
    sig_catch = ''
    interpret 'signal 'sav_catch
  end
  exit 'FAIL'

$get_service_pack: procedure
  call ApiSigErr 'ARG'

  hk = RegOpen('HLM', 'SOFTWARE\Microsoft\Windows NT\CurrentVersion')
  if hk = 'FAILED' then do
    say 'Unable to open registry LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion'
    return ''
  end

  ret = RegGetVal( hk,, 'CSDVersion', sp)
  if (ret ^= 'OK' | sp.1 = '') then do
    call RegClose ( hk )
    return ''
  end

  call RegClose ( hk )
  
  call ApiSigErr 'API'

  /* Handle correctly 'Dotatek Service Pack 3' and return simple 'Service Pack 3' */
  /* See bug 53596 */  
  service_pos = wordpos('Service', sp.1)
  if service_pos = 0 then return ''
  
  return subword(sp.1, service_pos)
  
$is_further_then_40:
  /* return 1 here is much safer */
  if system.version.major == 'SYSTEM.VERSION.MAJOR' | ,
     system.version.minor == 'SYSTEM.VERSION.MINOR' then return 1
  if system.version.major > 4 | ,
     (system.version.major == 4 & system.version.minor >= 10) then return 1
  return 0

$is_xp_or_further:
  /* return 0 here is much safer */
  if system.version.major == 'SYSTEM.VERSION.MAJOR' | ,
     system.version.minor == 'SYSTEM.VERSION.MINOR' then return 0
  if system.version.major > 5 | ,
     (system.version.major == 5 & system.version.minor >= 1) then return 1
  return 0
  
$is_vista_or_further:
  if system.version.major == 'SYSTEM.VERSION.MAJOR' | ,
     system.version.minor == 'SYSTEM.VERSION.MINOR' then return 0
  if system.version.major >= 6 then return 1
  return 0

$is_win2k_rollup_1_or_further:
  call ApiSigErr 'ARG'

  if system.os = '9X' then return 0
  if $is_xp_or_further() then return 1  
  if $is_further_then_40() then do
    rollup_1_hk = RegOpen('HLM', 'SOFTWARE\Microsoft\Updates\Windows 2000\SP5\Update Rollup 1')
    if rollup_1_hk <> 'FAILED' then return 1       
  end
  call ApiSigErr 'API'
   
  return 0

$is_spiderg3_possible:
  if system.version.major == 'SYSTEM.VERSION.MAJOR' | ,
     system.version.minor == 'SYSTEM.VERSION.MINOR' then do
    say '[Spider3G] system.version not defined!!! unable to determine spider3g possibility'
    return 0
  end

  if system.os = '9X' then do
    say '[Spider3G] spider3g is impossible on win9X'
    return 0
  end

  if $is_vista_or_further() then do
    say '[Spider3G] vista or further spider3g is possible'
    return 1
  end

  if system.version.major == 5 then do
    sp = $get_service_pack()
    
    if system.version.minor = 2 then do
      if (sp <> '') & (sp >= 'Service Pack 1') then do
        say '[Spider3G] 2k3 sp>1 spider3g is possible'
        return 1
      end
    end
    else if system.version.minor = 1 then do
      if (sp <> '') & (sp >> 'Service Pack 1') then do
        say '[Spider3G] XP sp>1 spider3g is possible'
        return 1
      end
    end
    else if $is_win2k_rollup_1_or_further() then do
      say '[Spider3G] 2k sp4 rollup1 spider3g is possible'
      return 1
    end   
  end
  say '[Spider3G] system is too ancient for spider3g. It is impossible'

  return 0
