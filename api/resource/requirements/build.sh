images_list=$1
rmp_list=$2

printf "开始处理镜像..."
for i in $(cat image_list);
do
  docker pull $i;
  docker tag localhost:8082/$i localhost:8082/$i;
  docker push localhost:8082/$i;
done
printf "镜像处理完毕!"

printf "开始处理rpm"
for r in $(cat rmp_list);
do
  yum install r -y
done
printf "rpm处理完毕！"

