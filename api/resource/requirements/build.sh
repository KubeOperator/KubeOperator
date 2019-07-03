registry="registry.fit2cloud.com/public"
printf "开始处理镜像..."
for i in $(cat $1);
do
  docker pull $registry/$i;
  docker tag  $registry/$i localhost:8092/$i;
  docker push localhost:8092/$i;
done
printf "镜像处理完毕!"

printf "开始处理rpm"
for r in $(cat $2);
do
  yumdownloader $r
done
printf "rpm处理完毕！"

