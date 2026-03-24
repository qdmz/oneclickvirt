#!/usr/bin/env python3
import mysql.connector
import json

# 旧数据库连接配置
old_db_config = {
    'host': 'localhost',
    'user': 'root',
    'database': 'oneclickvirt',
    'unix_socket': '/var/run/mysqld/mysqld.sock'
}

# 新数据库连接配置
new_db_config = {
    'host': 'localhost',
    'user': 'root',
    'database': 'oneclickvirt_new',
    'unix_socket': '/var/run/mysqld/mysqld.sock'
}

def get_connection(config):
    try:
        conn = mysql.connector.connect(**config)
        return conn
    except Exception as e:
        print(f"数据库连接失败: {e}")
        return None

def migrate_users(old_conn, new_conn):
    print("开始迁移用户数据...")
    cursor_old = old_conn.cursor(dictionary=True)
    cursor_new = new_conn.cursor()
    
    # 清空新数据库中的用户表
    cursor_new.execute("DELETE FROM users")
    
    # 从旧数据库获取用户数据
    cursor_old.execute("SELECT * FROM users")
    users = cursor_old.fetchall()
    
    # 迁移用户数据
    for user in users:
        try:
            # 构建INSERT语句
            sql = """
            INSERT INTO users (
                id, uuid, username, email, password, user_type, level, status, 
                nickname, phone, created_at, updated_at
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            """
            values = (
                user['id'], user['uuid'], user['username'], user['email'],
                user['password'], user['user_type'], user['level'], user['status'],
                user['nickname'], user['phone'], user['created_at'], user['updated_at']
            )
            cursor_new.execute(sql, values)
        except Exception as e:
            print(f"迁移用户 {user['username']} 失败: {e}")
    
    new_conn.commit()
    print(f"成功迁移 {len(users)} 个用户")

def migrate_providers(old_conn, new_conn):
    print("开始迁移节点数据...")
    cursor_old = old_conn.cursor(dictionary=True)
    cursor_new = new_conn.cursor()
    
    # 清空新数据库中的节点表
    cursor_new.execute("DELETE FROM providers")
    
    # 从旧数据库获取节点数据
    cursor_old.execute("SELECT * FROM providers")
    providers = cursor_old.fetchall()
    
    # 迁移节点数据
    for provider in providers:
        try:
            # 构建INSERT语句
            sql = """
            INSERT INTO providers (
                id, uuid, name, type, endpoint, ssh_port, username, password, 
                status, region, country, created_at, updated_at
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            """
            values = (
                provider['id'], provider['uuid'], provider['name'], provider['type'],
                provider['endpoint'], provider['ssh_port'], provider['username'], provider['password'],
                provider['status'], provider['region'], provider['country'],
                provider['created_at'], provider['updated_at']
            )
            cursor_new.execute(sql, values)
        except Exception as e:
            print(f"迁移节点 {provider['name']} 失败: {e}")
    
    new_conn.commit()
    print(f"成功迁移 {len(providers)} 个节点")

def migrate_instances(old_conn, new_conn):
    print("开始迁移实例数据...")
    cursor_old = old_conn.cursor(dictionary=True)
    cursor_new = new_conn.cursor()
    
    # 清空新数据库中的实例表
    cursor_new.execute("DELETE FROM instances")
    
    # 从旧数据库获取实例数据
    cursor_old.execute("SELECT * FROM instances")
    instances = cursor_old.fetchall()
    
    # 迁移实例数据
    for instance in instances:
        try:
            # 构建INSERT语句
            sql = """
            INSERT INTO instances (
                id, uuid, name, provider, provider_id, status, image, instance_type, 
                cpu, memory, disk, bandwidth, private_ip, public_ip, ssh_port, 
                username, password, os_type, expired_at, user_id, created_at, updated_at
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            """
            values = (
                instance['id'], instance['uuid'], instance['name'], instance['provider'],
                instance['provider_id'], instance['status'], instance['image'], instance['instance_type'],
                instance['cpu'], instance['memory'], instance['disk'], instance['bandwidth'],
                instance['private_ip'], instance['public_ip'], instance['ssh_port'],
                instance['username'], instance['password'], instance['os_type'],
                instance['expired_at'], instance['user_id'], instance['created_at'], instance['updated_at']
            )
            cursor_new.execute(sql, values)
        except Exception as e:
            print(f"迁移实例 {instance['name']} 失败: {e}")
    
    new_conn.commit()
    print(f"成功迁移 {len(instances)} 个实例")

def migrate_ports(old_conn, new_conn):
    print("开始迁移端口数据...")
    cursor_old = old_conn.cursor(dictionary=True)
    cursor_new = new_conn.cursor()
    
    # 清空新数据库中的端口表
    cursor_new.execute("DELETE FROM ports")
    
    # 从旧数据库获取端口数据
    cursor_old.execute("SELECT * FROM ports")
    ports = cursor_old.fetchall()
    
    # 迁移端口数据
    for port in ports:
        try:
            # 构建INSERT语句
            sql = """
            INSERT INTO ports (
                id, instance_id, provider_id, host_port, guest_port, protocol, 
                status, description, is_ssh, created_at, updated_at
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            """
            values = (
                port['id'], port['instance_id'], port['provider_id'], port['host_port'],
                port['guest_port'], port['protocol'], port['status'], port['description'],
                port['is_ssh'], port['created_at'], port['updated_at']
            )
            cursor_new.execute(sql, values)
        except Exception as e:
            print(f"迁移端口 {port['id']} 失败: {e}")
    
    new_conn.commit()
    print(f"成功迁移 {len(ports)} 个端口")

def migrate_orders(old_conn, new_conn):
    print("开始迁移订单数据...")
    cursor_old = old_conn.cursor(dictionary=True)
    cursor_new = new_conn.cursor()
    
    # 清空新数据库中的订单表
    cursor_new.execute("DELETE FROM orders")
    
    # 从旧数据库获取订单数据
    cursor_old.execute("SELECT * FROM orders")
    orders = cursor_old.fetchall()
    
    # 迁移订单数据
    for order in orders:
        try:
            # 构建INSERT语句
            sql = """
            INSERT INTO orders (
                id, uuid, user_id, product_id, amount, status, 
                payment_method, payment_time, created_at, updated_at
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            """
            values = (
                order['id'], order['uuid'], order['user_id'], order['product_id'],
                order['amount'], order['status'], order['payment_method'], order['payment_time'],
                order['created_at'], order['updated_at']
            )
            cursor_new.execute(sql, values)
        except Exception as e:
            print(f"迁移订单 {order['id']} 失败: {e}")
    
    new_conn.commit()
    print(f"成功迁移 {len(orders)} 个订单")

def migrate_user_wallets(old_conn, new_conn):
    print("开始迁移钱包数据...")
    cursor_old = old_conn.cursor(dictionary=True)
    cursor_new = new_conn.cursor()
    
    # 清空新数据库中的钱包表
    cursor_new.execute("DELETE FROM user_wallets")
    
    # 从旧数据库获取钱包数据
    cursor_old.execute("SELECT * FROM user_wallets")
    wallets = cursor_old.fetchall()
    
    # 迁移钱包数据
    for wallet in wallets:
        try:
            # 构建INSERT语句
            sql = """
            INSERT INTO user_wallets (
                id, user_id, balance, created_at, updated_at
            ) VALUES (%s, %s, %s, %s, %s)
            """
            values = (
                wallet['id'], wallet['user_id'], wallet['balance'],
                wallet['created_at'], wallet['updated_at']
            )
            cursor_new.execute(sql, values)
        except Exception as e:
            print(f"迁移钱包 {wallet['id']} 失败: {e}")
    
    new_conn.commit()
    print(f"成功迁移 {len(wallets)} 个钱包")

def main():
    # 连接旧数据库
    old_conn = get_connection(old_db_config)
    if not old_conn:
        print("无法连接旧数据库")
        return
    
    # 连接新数据库
    new_conn = get_connection(new_db_config)
    if not new_conn:
        print("无法连接新数据库")
        old_conn.close()
        return
    
    try:
        # 迁移数据
        migrate_users(old_conn, new_conn)
        migrate_providers(old_conn, new_conn)
        migrate_instances(old_conn, new_conn)
        migrate_ports(old_conn, new_conn)
        migrate_orders(old_conn, new_conn)
        migrate_user_wallets(old_conn, new_conn)
        
        print("数据迁移完成!")
    except Exception as e:
        print(f"迁移过程中出错: {e}")
    finally:
        # 关闭数据库连接
        old_conn.close()
        new_conn.close()

if __name__ == "__main__":
    main()
