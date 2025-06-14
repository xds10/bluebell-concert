import mysql.connector
from mysql.connector import Error

# ========== 配置参数 ==========
# 数据库连接配置
DB_CONFIG = {
    'host': 'localhost',
    'user': 'root',
    'password': '123456',
    'database': 'concert-ticket-sales',
    'port': 8086,
    'charset': 'utf8mb4'
}

# 场馆ID
PLACE_ID = 2

PLACE="大型体育场"

# 座位区域配置（区域名称列表）
SECTIONS = ['VIP区', 'A区', 'B区', 'C区']

# 每个区域的排数
ROWS_PER_SECTION = 50

# 每排的座位数
SEATS_PER_ROW = 100

# 区域价格配置（需与区域数量匹配）
PRICES = [1980.0, 1280.0, 880.0, 580.0]

# 数据库表字段名映射
FIELD_MAPPING = {
    'section': 'section',     # 区域字段名
    'seat_row': 'seat_row',             # 排号字段名
    'seat_no': 'seat_no',     # 座位号字段名
    'price': 'price',         # 价格字段名
}
# ==============================

def generate_seats():
    """生成座位数据"""
    seats = []
    for section_idx, section in enumerate(SECTIONS):
        price = PRICES[section_idx]
        for row in range(1, ROWS_PER_SECTION + 1):
            for seat_no in range(1, SEATS_PER_ROW + 1):
                seats.append({
                    'place_id': PLACE_ID,
                    'place':PLACE,
                    'section': section,
                    'seat_row': row,
                    'seat_no': seat_no,
                    'price': price
                })
    return seats

def insert_to_mysql(seats):
    """将座位数据插入MySQL数据库"""
    try:
        connection = mysql.connector.connect(**DB_CONFIG)
        if connection.is_connected():
            cursor = connection.cursor()
            
            # 构建插入SQL
            field_names = ', '.join(['place_id']+['place'] + list(FIELD_MAPPING.values()))
            placeholders = ', '.join(['%s'] * (2 + len(FIELD_MAPPING)))
            insert_query = f"INSERT INTO seat ({field_names}) VALUES ({placeholders})"
            
            # 批量插入
            data_to_insert = []
            for seat in seats:
                row_data = [
                    seat['place_id'],
                    seat['place'],
                    seat['section'],
                    seat['seat_row'],
                    seat['seat_no'],
                    seat['price']
                ]
                data_to_insert.append(tuple(row_data))
            
            cursor.executemany(insert_query, data_to_insert)
            connection.commit()
            print(f"成功插入 {cursor.rowcount} 条座位数据")
            
    except Error as e:
        print(f"数据库错误: {e}")
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()
            print("数据库连接已关闭")

def main():
    # 生成座位数据
    seats = generate_seats()
    
    # 插入数据库
    insert_to_mysql(seats)

if __name__ == "__main__":
    main()    