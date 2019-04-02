#! /bin/bash

red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}

i=0

colorMsg $yellow  "Which package do you want to build?"

for file in ./*
do
    if test -d $file
    then
        let i++
	colorMsg $green "     $i)${file:2}"
    fi
done
let i++

flag=0
while [ "$flag"x == "0"x ]
do
read -p "Please input your option:" x
j=0
for file in ./*
do
    if test -d $file
    then
        let j++
	if [ "$x"x == "$j"x ]  
        then
            let flag++
             cd ./$file
            ./build.sh 
        fi
    fi
done
if [ "$flag"x == "0"x ]
then
	colorMsg $red "input error: $x"
fi
done

