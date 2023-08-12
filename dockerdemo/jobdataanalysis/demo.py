import pandas as pd;
df = pd.read_csv('/data/Passanger_booking_data.csv');
print(df.head())

# 路线分析：统计各个路线的预订数量
route_counts = df['route'].value_counts()

# 找出最受欢迎的路线
most_popular_route = route_counts.idxmax()
most_popular_count = route_counts.max()

# 在stdout上输出结果
print('路线预订数量统计：')
print(route_counts.head(10))  # 输出前10个路线的统计数据作为示例
print('\n最受欢迎的路线：')
print(f'路线：{most_popular_route}')
print(f'预订数量：{most_popular_count}')
