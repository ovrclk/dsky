
declare -a deps
declare -A dep_1
deps[1]=dep_1
dep_1[name]=foo
dep_1[state]=active


for dep in ${deps[@]}
do
  eval "echo \$${dep}[state]"
done
