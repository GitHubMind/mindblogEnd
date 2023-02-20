RED_COLOR='\E[1;31m'  #红
GREEN_COLOR='\E[1;32m' #绿
YELOW_COLOR='\E[1;33m' #黄
BLUE_COLOR='\E[1;34m'  #蓝
PINK='\E[1;35m'      #粉红
RES='\E[0m'

echo -e  "${GREEN_COLOR}****开始执行自动化部署****${RES}\n\n"

echo -e "${YELOW_COLOR}---step1:合并代码---${RES}"

git pull

echo -e "${BLUE_COLOR}合并代码成功${RES}\n"

echo -e "${YELOW_COLOR}---step2:编译---${RES}"

go build -o blog main.go

echo -e "${BLUE_COLOR}编译完成${RES}\n"

echo -e "${YELOW_COLOR}---step3:更改权限---${RES}"
chmod -R 777 blog
echo -e "${BLUE_COLOR}更改权限完成${RES}\n"


i1=`ps -ef|grep -E "blog"|grep -v grep|awk '{print $2}'`
if [ $i1 ]
then
echo -e "${YELOW_COLOR}---step4:杀掉进程并且运行---${RES}"
echo -e "${BLUE_COLOR}杀掉进程$i1${RES}\n"
kill -9 $i1
fi
#   9       KILL (non-catchable, non-ignorable kill)



nohup ./blog >bloglogs 2>&1 &

i2=`ps -ef|grep -E "blog"|grep -v grep|awk '{print $2}'`

echo -e "${GREEN_COLOR}****部署成功,部署的进程ID为:$i2${RES}****"