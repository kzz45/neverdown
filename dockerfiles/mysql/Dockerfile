FROM mysql:5.7

LABEL maintainer="kongzz"

RUN yum -y install hostname vim

COPY mysqld.cnf /etc/mysql/mysql.conf.d/mysqld.cnf

COPY run.sh /

CMD ["bash", "-x", "/run.sh"]
