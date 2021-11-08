COUNTER=$1
for i in `seq 0 $COUNTER`
do
  port=$(($i+1234))
  lsof -t -i tcp:$port | xargs kill -s
done