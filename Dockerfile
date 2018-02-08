# 使用1.9的golang作为母镜像
FROM golang:1.9

# 维护者信息，也就是作者的姓名和邮箱
MAINTAINER eyeblue "eyebluecn@126.com"

# 指定工作目录就是 compass。工作目录指的就是以后 每层构建的当前目录 。
WORKDIR $GOPATH/src/compass

# 将compass项目下的所有文件移动到golang镜像中去
COPY . $GOPATH/src/compass

# 这里为了维持docker无状态性，准备数据卷作为日志目录和上传文件目录
VOLUME /data/log
VOLUME /data/matter
# 通过环境变量的方式，为应用指定日志目录和上传文件目录。
ENV COMPASS_LOG_PATH=/data/log COMPASS_MATTER_PATH=/data/matter


# golang.org库国内无法下载，这里从我准备的github中clone
# github.com的库可以直接通过`go get`命令下载
# `go install compass`是对项目进行打包
# `cp`是将项目需要的html等文件移动到可执行文件的目录下。
RUN git clone https://github.com/eyebluecn/golang.org.git $GOPATH/src/golang.org \
    && go get github.com/disintegration/imaging \
    && go get github.com/json-iterator/go \
    && go get github.com/go-sql-driver/mysql \
    && go get github.com/jinzhu/gorm \
    && go get github.com/nu7hatch/gouuid \
    && go install compass \
    && cp -r $GOPATH/src/compass/build/* $GOPATH/bin

# 声明运行时容器提供服务端口，这只是一个声明。默认是6030端口
EXPOSE 6030

# compass作为执行文件 启动这个容器就会去执行 `/go/bin/compass`
ENTRYPOINT ["/go/bin/compass"]
