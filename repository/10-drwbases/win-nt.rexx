/* drwbases: win-nt.rexx */

signal on syntax  name log_sig
signal on error   name log_sig
signal on failure name log_sig
call ApiSigErr 'API'

/* ------------------------------------------------------------------ */

exit $do_upgrade()

/* ------------------------------------------------------------------ */

$instdir: procedure
  return $instdir_home( arg( 1 ) )

$before_files:
  return $before_files_win( arg( 1 ) )

$after_files:
  return $after_files_win( arg( 1 ) )
