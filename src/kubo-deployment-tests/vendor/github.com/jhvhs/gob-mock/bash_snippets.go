package gobmock

const (
	scriptStart    = "%[1]s() {\n"
	stubDefinition = `
          while read -t0.05; do
            :
          done
  	`
	spyDefinition = `
          local in_line_count=0
          declare -a in_lines
          while read -r -t0.1; do
            in_lines[in_line_count]="$REPLY"
            in_line_count=$((in_line_count + 1))
          done
          callCounter=$((callCounter + 1))
          echo "<$callCounter> %[1]s $*" > /dev/fd/2
          if [ ${in_line_count} -gt 0 ]; then
            echo "<$callCounter received> input:" > /dev/fd/2
            printf '%%s\n' "${in_lines[@]}" > /dev/fd/2
            echo "<$callCounter end received>" > /dev/fd/2
          fi
        `
	mockDefinition   = "%[2]s\n"
	scriptEnd        = "\n return 0 \n } \n"
	exportDefinition = "export -f %s\n"

	callThroughDefinition = `
	  if [ ${in_line_count} -gt 0 ]; then
	    echo $(IFS=$'\n'; echo ${in_lines[*]}) | $(which %[1]s) "$@"
	  else
	    $(which %[1]s) "$@"
	  fi
  	`
)
