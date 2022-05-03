FROM jrei/systemd-centos:7

COPY OracleManager.viper.yaml /etc/OracleManager/OracleManager.viper.yaml
COPY oracle-manager /usr/local/bin/oracle-manager

CMD ["oracle-manager", "run"]
