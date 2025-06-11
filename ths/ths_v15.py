'''
exe:
  pip install pyinstaller 
  
  pyinstaller --onefile --clean --add-data 'cfg.toml;.' --name ths_cmd.exe ths_v15.py 
  pyinstaller --onefile --clean --add-data 'cfg.toml;.' --name ths_nogui.exe --noconsole ths_v15.py 

'''

import os
import pandas as pd

# from validators import domain
import datetime,requests,time
from io import BytesIO
from retry import retry 
from chinese_calendar import is_holiday
import multitasking

import toml

#=====================================================================
INDEX_LIST = ['SHSE.000001','SHSE.000300','SZSE.399001','SZSE.399006','SHSE.510880']
#=====================================================================

ss_ths_1d = '''

from datetime import datetime,timedelta

V1_LIST   = []
OO_LIST   = []
CC_LIST   = []

for i in range(0,total):
    proo=get("OPEN", i)       
    prcc=get("CLOSE", i)      
    V0=get("VOLUME",i)/100    
        
    nt=int(get("TIME", i))  # eg. 20240412
    tkey1 = str(nt)
    
    if tkey1 in barsdata:
        row  = barsdata[tkey1]
        proo=row[0]  
        prcc=row[1] 
        vv1 =row[2]  
        
        V1_LIST.append(vv1)
        OO_LIST.append(proo)
        CC_LIST.append(prcc)
        continue
    else: 
        dt = datetime.fromtimestamp(nt)#-timedelta(hours=8)
        tkey = str(dt.strftime('%Y%m%d %H:%M:%S'))
        
        if tkey in bars1m:
            row  = bars1m[tkey]
            proo=row[0]  
            prcc=row[1] 
            vv1 =row[2]  
            
            V1_LIST.append(vv1)
            OO_LIST.append(proo)
            CC_LIST.append(prcc)
            
        else: 
            V1_LIST.append(0)
            OO_LIST.append(0)
            CC_LIST.append(0)


for i in range(1,total):
    SL0 = V1_LIST[i]
    SL1 = V1_LIST[i-1]
    hevo.save("V1", SL0, i) 
    hevo.save("V2", SL0+SL1, i) 
    
    if OO_LIST[i]<CC_LIST[i]:
        hevo.save("UP", SL0, i) 
    else:
        hevo.save("DOWN", SL0, i)         
        
draw.stick("V2",14, 1)
draw.stick("V1", 5, 1)
# draw.curve_right("PE","#735595")

draw.stick("UP", 4, 1)
draw.stick("DOWN", 8, 2)

'''

ss_ths_1m = '''

#=====================
OO_LIST   = []
CC_LIST   = []
VV_LIST   = []

V1_LIST   = []
V2_LIST   = []
V3_LIST   = []

# num = param(5)

def calculate_median(nums):
    if not nums:
        raise ValueError("The list is empty. Median cannot be calculated.")
    
    # 排序列表
    nums_sorted = sorted(nums)
    n = len(nums_sorted)
    
    # 奇数长度
    if n % 2 == 1:
        return nums_sorted[n // 2]
    # 偶数长度
    else:
        mid1 = nums_sorted[n // 2 - 1]
        mid2 = nums_sorted[n // 2]
        return (mid1 + mid2) / 2
    
#=====================
for i in range(0,total):
    oo = get("OPEN", i)
    cc = get("CLOSE", i)
    vv = get("VOLUME",i)/100
    
    ndate=int(get("TIME", i))  # eg. 20240412
    tkey = str(ndate)
    OO_LIST.append(oo)
    CC_LIST.append(cc)
    VV_LIST.append(vv)
    
    if tkey in barsdata:
        data  = barsdata[tkey]
        v1 = data[0]
        v2 = data[1]
        v3 = data[2]
        # pet  = data[4]
        
        V1_LIST.append(v1)
        V2_LIST.append(v1+v2)
        V3_LIST.append(v3)
        
    else: 
        V1_LIST.append(0)
        V2_LIST.append(0)
        V3_LIST.append(0)

#=====================
for i in range(1,total):
    oo  = OO_LIST[i]
    cc  = CC_LIST[i]
    v11 = V1_LIST[i]
    v10 = V1_LIST[i-1]
    v22 = V2_LIST[i]
    v33 = V3_LIST[i]
    
    hevo.save("V1", v11, i) 
    hevo.save("V2", v22, i) 
    hevo.save("V3", v33, i) 
    
    if i>1 and v10>0:
        hevo.save("VV%", 100.0*v11/v10, i) 
        
    if v11>0:
        hevo.save("V93%", 100.0*v33/v11, i) 
    
    if oo<cc:
        hevo.save("U", v11, i) 
    else:
        hevo.save("D", v11, i) 
        
    # Show number of ratio for v931/v931[-1]
    if v11>v10*1.99 and v10>0:
        msg = str(round(v11/v10,1))
        text(V1_LIST[i]*1.3, i, msg, 3) # To set ratio value above the bars 

# To show a horizontal line for v150 
for i in range(begin,end): 
    v3 = V3_LIST[i] 
    draw.line(v3, i, v3, i+1, "#FF000D") 
    
#=====================
draw.stick("V2",14, 1)
draw.stick("V1", 5, 1)

# draw.color_stick("V3","#735595",2) 
draw.color_stick("V3") 
# draw.curve("V3") 
draw.curve_right("VV%", 9, 0)  # set to no-draw 
draw.curve_right("V93%", 5, 0) # set to no-draw 
# draw.curve_right("PE","#735595")

draw.stick("D", 8, 2)
draw.stick("U", 4, 1)



'''

ss_ths_5m = '''

V1_LIST   = []
V3_LIST   = []
# 成交量列表 = []

for i in range(0,total):
    proo=get("OPEN", i)       #获取每条K线上的开盘价
    prcc=get("CLOSE", i)      #获取每条K线上的收盘价
    V0=get("VOLUME",i)/100      # 当前K线成交量
        
    ndate=int(get("TIME", i))  # eg. 20240412
    tkey = str(ndate)
    
    if tkey in barsdata:
        row  = barsdata[tkey]
        vv1 = row[0]
        vv2 = row[1]
        pet = row[2]
        
        V1_LIST.append(vv1)
        V3_LIST.append(vv2)
        
        hevo.save("V1", vv1, i) 
        hevo.save("V2", vv1+vv2, i) 
            
        hevo.save("PE", pet, i) 
        
        if proo<prcc:
            hevo.save("UP", vv1, i) 
        else:
            hevo.save("DOWN", vv1, i) 
    else: 
        V1_LIST.append(0)
        V3_LIST.append(0)

for i in range(1,total):
    SL0 = V1_LIST[i]
    SL1 = V1_LIST[i-1]
        
    if SL0>SL1*1.99 and SL1>0:
        # 成交量与首量背离，且首量缩量，而日成交量放量：此时表示主力没有在早盘出手，且盘中有大卖盘，所以应卖出
        msg = str(round(SL0/SL1,1))
        text(V1_LIST[i]*1.3, i, msg, 3)

# for i in range(begin,end):    
#     SLL = V3_LIST[i]    
#     draw.line(SLL, i, SLL, i+1, "#FF000D")  #画一条数值为尾量的直线到下一个K线
        
draw.stick("V2",14, 1)
draw.stick("V1", 5, 1)
draw.curve_right("PE","#735595")

draw.stick("UP", 4, 1)
draw.stick("DOWN", 8, 2)

'''

ss_ths_vr = '''

#=====================
OO_LIST   = []
CC_LIST   = []
VV_LIST   = []

V1_LIST   = []
V2_LIST   = []
V3_LIST   = []

num = param(5)

def calculate_median(nums):
    if not nums:
        raise ValueError("The list is empty. Median cannot be calculated.")
    
    nums_sorted = sorted(nums)
    n = len(nums_sorted)
    
    if n % 2 == 1:
        return nums_sorted[n // 2]
    else:
        mid1 = nums_sorted[n // 2 - 1]
        mid2 = nums_sorted[n // 2]
        return (mid1 + mid2) / 2
    
#=====================
for i in range(0,total):
    oo = get("OPEN", i)
    cc = get("CLOSE", i)
    vv = get("VOLUME",i)/100
    
    ndate=int(get("TIME", i))  # eg. 20240412
    tkey = str(ndate)
    OO_LIST.append(oo)
    CC_LIST.append(cc)
    VV_LIST.append(vv)
    
    if tkey in barsdata:
        data  = barsdata[tkey]
        v1 = data[0]/100
        v2 = data[1]/100
        v3 = data[2]/100
        # pet  = data[4]
        
        V1_LIST.append(v1)
        V2_LIST.append(v2)
        V3_LIST.append(v3)
        
    else: 
        V1_LIST.append(0)
        V2_LIST.append(0)
        V3_LIST.append(0)

vmed_list = VV_LIST[:num]
#=====================
for i in range(num,total):
    oo  = OO_LIST[i]
    cc  = CC_LIST[i]
    vv  = VV_LIST[i]
    
    v11 = V1_LIST[i]
    v10 = V1_LIST[i-1]
    v22 = V2_LIST[i]
    v33 = V3_LIST[i]
    
    vmed = calculate_median(VV_LIST[i-num:i])
    
    # hevo.save("V1", v11, i) 
    # hevo.save("V2", v22, i) 
    # hevo.save("V3", v33, i) 
    
    rr1 = 100.0*v11/vmed 
    rr2 = 100.0*(v11+v22)/vmed 
    rr3 = 100.0*v33/vmed 
    
    vmed_list.append(rr1)
    
    if vmed>0:
        hevo.save("V1/Vmed", rr1, i) 
        hevo.save("V2/Vmed", rr2, i) 
        hevo.save("V3/Vmed", rr3, i) 

    # Note: This should be the last save params. 
    if oo<cc:
        hevo.save("U", rr1, i) 
    else:
        hevo.save("D", rr1, i) 
    
#=====================
vr_max = max(vmed_list[begin:end])
if vr_max>50:
    draw.line(50, begin, 50, end, "#507efbb3") 

# if vr_max>30:
#     draw.line(30, begin, 30, end, "#507efbb3") 

if vr_max>20:
    draw.line(20, begin, 20, end, "#507efbb3") 

if vr_max>10:
    draw.line(10, begin, 10, end, "#507efbb3") 

# if vr_max>5:
#     draw.line(5,  begin, 5,  end, "#507efbb3") 

draw.curve_right("V3/Vmed", "#80ef1de7", 1)  # set to no-draw 
# draw.curve("V1/v", "#30c9ff27", 1)  # set to no-draw 
# draw.curve_right("PE","#735595")

draw.stick("V2/Vmed",14, 1)
draw.stick("V1/Vmed", 5, 1)

# draw.color_stick("V3","#735595",2) 
# draw.color_stick("V3") 
# draw.curve("V3") 

draw.stick("D", 8, 2)
draw.stick("U", 4, 1)

'''

ss_ths_ud = '''

#=====================
V1_LIST   = []
V2_LIST   = []
# 成交量列表 = []

#=====================
for i in range(0,total):
    proo=get("OPEN", i)       #获取每条K线上的开盘价
    prcc=get("CLOSE", i)      #获取每条K线上的收盘价
    V0=get("VOLUME",i)/100      # 当前K线成交量
    
    # 成交量列表.append(V0)
    
    ndate=int(get("TIME", i))  # eg. 20240412
    tkey = str(ndate)
    
    if tkey in barsdata:
        row  = barsdata[tkey]
        vv1 = row[0]
        vv2 = row[1]
        
        V1_LIST.append(vv1)
        V2_LIST.append(vv2)
        
        hevo.save("U1", vv1, i) 
        hevo.save("U2", vv2, i) 
        hevo.save("A2", vv1+vv2, i) 
            
        if proo<prcc:
            hevo.save("U", vv1, i) 
        else:
            hevo.save("D", vv1, i) 
    else: 
        V1_LIST.append(0)
        V2_LIST.append(0)

#=====================
for i in range(1,total):
    v11 = V1_LIST[i]
    v10 = V1_LIST[i-1]
        
    if v11>v10*1.99 and v10>0:
        msg = str(round(v11/v10,1))
        text(V1_LIST[i]*1.3, i, msg, 3)

# for i in range(begin,end):    
#     SLL = V3_LIST[i]    
#     draw.line(SLL, i, SLL, i+1, "#FF000D")  #

#=====================
draw.stick("A2",14, 1)
draw.stick("U1", 5, 1)

draw.curve("U2", 6, 1)

draw.stick("U", 4, 1)
draw.stick("D", 8, 2)


'''


ss_ths_pvj = '''

V1_LIST   = []
for i in range(0,total):
    proo=get("OPEN", i)       
    prcc=get("CLOSE", i)      
    V0=get("VOLUME",i)/100 
    
    ndate=int(get("TIME", i))  # eg. 20240412
    tkey = str(ndate)
    
    if tkey in barsdata:
        row  = barsdata[tkey]
        vv1 = row[0]
        
        V1_LIST.append(vv1)        
        hevo.save("PVJ", vv1, i) 
        
    else: 
        V1_LIST.append(0)

draw.curve("PVJ")

'''

ss_ths_cbj = '''

#=====================
OO_LIST   = []
CC_LIST   = []
# VV_LIST   = []

V1_LIST   = []
V2_LIST   = []
V3_LIST   = []


#=====================
for i in range(0,total):
    oo = get("OPEN", i)
    cc = get("CLOSE", i)
    # vv = get("VOLUME",i)/100
    
    ndate=int(get("TIME", i))  # eg. 20240412
    tkey = str(ndate)
    OO_LIST.append(oo)
    CC_LIST.append(cc)
    # VV_LIST.append(vv)
    
    if tkey in barsdata:
        data  = barsdata[tkey]
        c0 = data[0]
        c1 = data[1]
        c2 = data[2]
        # pet  = data[4]
        
        # V1_LIST.append(c0)
        # V2_LIST.append(c1)
        # V3_LIST.append(c2)
        
        save("CBJ", c0, i) 
        save("CBh", c1, i) 
        save("CBt", c2, i) 
        
    # else: 
    #     V1_LIST.append(0)
    #     V2_LIST.append(0)
    #     V3_LIST.append(0)
    
draw.curve("CBh", "#7500ffff")
draw.curve("CBJ", "#95ffff14")
draw.curve("CBt", "#75ff81c0")

'''

ss_ths_cbf = '''

#=====================
# OO_LIST   = []
# CC_LIST   = []
# VV_LIST   = []

V1_LIST   = []
V2_LIST   = []
V3_LIST   = []


#=====================
for i in range(0,total):
    # oo = get("OPEN", i)
    # cc = get("CLOSE", i)
    # vv = get("VOLUME",i)/100
    
    ndate=int(get("TIME", i))  # eg. 20240412
    tkey = str(ndate)
    # OO_LIST.append(oo)
    # CC_LIST.append(cc)
    # VV_LIST.append(vv)
    
    if tkey in barsdata:
        data  = barsdata[tkey]
        c0 = data[0]
        c1 = data[1]
        c2 = data[2]
        # pet  = data[4]
        
        # V1_LIST.append(c0)
        # V2_LIST.append(c1)
        # V3_LIST.append(c2)
        
        save("c-t", c0-c2, i) 
        save("c-h", c0-c1, i) 
        
        ff = abs(c2-c1)
        if ff>0:
            save("hrr(%)", 100.0*abs(c0-c1)/ff, i) 
        else:
            save("hrr(%)", 0, i) 
        
    # else: 
    #     V1_LIST.append(0)
    #     V2_LIST.append(0)
    #     V3_LIST.append(0)
    
# draw.curve("CBh", "#7500ffff")
# draw.curve("CBJ", "#95ffff14")
draw.curve_right("hrr(%)", "#75ff81c0",0)

draw.stick("c-t", 8,2) 
draw.stick("c-h", 0,2) 

# draw.color_stick("h-c") 
# draw.color_stick("t-c") 
'''


to_drop_utc = lambda tt_utc: pd.to_datetime(tt_utc).tz_localize(None)
''' 去除UTC信息 '''


def get_trading_days(start_date=None, end_date=None):
    '''返回交易日历'''
    now = datetime.datetime.now()
    if start_date is None: start_date = datetime.date(now.year,1,1)
    if end_date is None: 
        end_date = datetime.date(now.year,12,31)
    elif not isinstance(end_date,datetime.date):
        end_date = pd.to_datetime(end_date).to_pydatetime().date()
    if not isinstance(start_date,datetime.date):
        start_date = pd.to_datetime(start_date).to_pydatetime().date()
    
    current_date = start_date
    trading_days = []
    while current_date <= end_date:
        # 排除周末和节假日
        if not is_holiday(current_date) and current_date.weekday() < 5:
            trading_days.append(current_date)
        current_date += datetime.timedelta(days=1)
        
    if len(trading_days)>0: trading_days = [str(day) for day in trading_days]
    
    return trading_days

def get_not_trading_days(start_date=None, end_date=None):
    '''返回交易日历'''
    now = datetime.datetime.now()
    if start_date is None: start_date = datetime.date(now.year,1,1)
    if end_date is None: 
        end_date = datetime.date(now.year,12,31)
    elif not isinstance(end_date,datetime.date):
        end_date = pd.to_datetime(end_date).to_pydatetime().date()
    if not isinstance(start_date,datetime.date):
        start_date = pd.to_datetime(start_date).to_pydatetime().date()
    
    current_date = start_date
    trading_days = []
    while current_date <= end_date:
        # 排除周末和节假日
        if is_holiday(current_date) or current_date.weekday() >= 5:
            trading_days.append(current_date)
        current_date += datetime.timedelta(days=1)
        
    if len(trading_days)>0: trading_days = [str(day) for day in trading_days]
    
    return trading_days


def get_trade_list_by_year(year=None):
    ''' 返回特定年份内的交易日列表 '''
    if year is None: year = pd.Timestamp.now().year
    
    sdate = f'{year}-01-01'
    edate = pd.to_datetime(sdate)+pd.tseries.offsets.YearEnd()
    edate = str(edate.date())
    # print(f' {sdate} ~ {edate}')
    trade_dates_list = get_trading_days(sdate,edate)
    return trade_dates_list

def get_trade_list_by_month(year=None,month=None):
    ''' 返回特定月份内的交易日列表 '''
    if year is None: year = pd.Timestamp.now().year
    if month is None: month = pd.Timestamp.now().month
    
    sdate = f'{year}-{month:02d}-01'
    edate = pd.to_datetime(sdate)+pd.tseries.offsets.MonthEnd()
    edate = str(edate.date())
    # print(f' {sdate} ~ {edate}')
    trade_dates_list = get_trading_days(sdate,edate)
    return trade_dates_list


def get_url_without_retry(url,params,connect_timeout=None,data_timeout=None):
    """
    不带重试的请求函数。
    如果请求失败或状态码异常，将抛出异常。
    
    参数：
        connect_timeout: 连接超时时间(默认3s)
        data_timeout: 传输超时时间(默认=connect_timeout)
    """
    if data_timeout is None: data_timeout = connect_timeout 
    response = requests.get(url,params=params, timeout=(connect_timeout,data_timeout))  # 设置超时时间
    response.raise_for_status()  # 如果状态码非200，抛出HTTPError
    return response.json()

# 定义重试机制：最多重试5次，间隔3秒
@retry(tries=5, delay=3, backoff=1, exceptions=(requests.exceptions.RequestException,))
def get_url_with_retry(url,params,connect_timeout=None,data_timeout=None):
    """
    带重试的请求函数：最多重试5次，重试间隔3秒。
    如果请求失败或状态码异常，将抛出异常。
    
    参数：
        connect_timeout: 连接超时时间(默认3s)
        data_timeout: 传输超时时间(默认=connect_timeout)
    """
    if data_timeout is None: data_timeout = connect_timeout 
    response = requests.get(url,params=params, timeout=(connect_timeout,data_timeout))  # 设置超时时间
    response.raise_for_status()  # 如果状态码非200，抛出HTTPError
    return response.json()


def get_date_interval_list(sdate=None,edate=None,freq='ME',name='date',include='both',normalize=True,): 
    '''
    - 获取日期区间列表
    - include: [both,left,right]
    '''
    if sdate is None: sdate = '2005-01-01'
    if edate is None: edate = str(pd.Timestamp.today().date())
    
    drng = pd.date_range(sdate,edate,freq=freq,inclusive='both',normalize=normalize,name=name)

    dtps = pd.Series(drng)
    
    if include!='right' and pd.to_datetime(sdate) not in dtps.values : 
        dtps.loc[len(dtps)] = pd.Timestamp(sdate)
        
    if include!='left'  and pd.to_datetime(edate) not in dtps.values : 
        dtps.loc[len(dtps)] = pd.Timestamp(edate)
        
    dtps = dtps.sort_values(ignore_index=True)
    
    return dtps 



def fetch_gm__calendar(eyear=None,syear=None,exchange=None,parse_dates=False,
                    ip='localhost',port=5000,timeout=5,data_timeout=None,retry=True,debug=False):
    '''
    请求交易日历(gm-api)
    
    参数：
        eyear: 结束年份，默认为当年
        syear: 开始年份, 默认为2005
        timeout: 连接超时时间(默认5s)
        data_timeout: 连接传输超时时间(默认timeout*2)
    
    '''
    if data_timeout is None: data_timeout = timeout*2

    if port==443: url = f'https://{ip}/get_dates_by_year'
    else: url = f'http://{ip}:{port}/get_dates_by_year'
    
    if eyear is None: eyear =  datetime.datetime.now().year

    params = {
        'eyear': str(eyear)
        }
    if syear is not None: 
        if syear>eyear: syear = eyear 
        params = dict(syear=str(syear),**params)
    if exchange is not None: params = dict(exchange=exchange,**params)

    try:
        if retry: json_data = get_url_with_retry(url,params,timeout,data_timeout=data_timeout)
        else: json_data = get_url_without_retry(url,params,timeout,data_timeout=data_timeout)
        df = pd.DataFrame(json_data)
        if df.shape[0]>0:
            df = df.set_index('date')
            df.index = pd.to_datetime(df.index)
            if parse_dates: df = df.apply(pd.to_datetime)
            df['is_trading'] = (~pd.to_datetime(df.trade_date).isna()).astype(int)
            # elif dtype is not None: df = df.astype(object).astype(str)            
        return df
        
    except (requests.exceptions.RequestException, ValueError, KeyError) as e:
        if debug: print(f" !!! Exception: {e}")
        return pd.DataFrame()
    


def fetch_gm__market_infos(sec='stock',exchange=None,onlytrade=True,filt_fields=True,
                    ip='localhost',port=5000,timeout=5,data_timeout=None,retry=True,debug=False):
    '''
    请求市场个股列表(gm-api)
    
    参数：
        sec: 市场类型，[stock,fund,index,...]
        timeout: 连接超时时间(默认5s)
        data_timeout: 连接传输超时时间(默认timeout*2)    
    '''
    if data_timeout is None: data_timeout = timeout*2

    if port==443: url = f'https://{ip}/get_infos'
    else: url = f'http://{ip}:{port}/get_infos'

    params = {}
    if sec is not None: params = dict(sec=sec,**params)
    if exchange is not None: params = dict(exchange=exchange,**params)
    
    try:
        if retry: json_data = get_url_with_retry(url,params,timeout,data_timeout=data_timeout)
        else: json_data = get_url_without_retry(url,params,timeout,data_timeout=data_timeout)
        
        df = pd.DataFrame(json_data)
        if df.shape[0]>0:
            df.listed_date = df.listed_date.apply(to_drop_utc)
            df.delisted_date = df.delisted_date.apply(to_drop_utc)
            
            if onlytrade : 
                today = datetime.datetime.now().strftime('%Y-%m-%d')
                df = df.loc[df['delisted_date']>today]
                
            if 'symbol' in df.columns: df.set_index('symbol',inplace=True)
            if filt_fields:
                fields = ['exchange','sec_id','sec_name','sec_abbr','sec_type1','sec_type2',
                        'listed_date','delisted_date']
                df = df[fields]
        else:
            if debug: print(f" !!! Warning: blank DataFrame for {url} - {sec}")
            
        return df
        
    except (requests.exceptions.RequestException, ValueError, KeyError) as e:
        if debug: print(f" !!! Exception: {e}")
        return pd.DataFrame()

def fetch_gm__daily_valuation(symbol,sdate=None,edate=None,inclusive='both',
                indicators='pe_ttm,pe_lyr,pb_lyr,ps_ttm,ps_lyr,dy_ttm,dy_lfy',
                 ip='localhost',port=5000,timeout=10,data_timeout=None,retry=True,debug=False):
    '''
    个股估值指标每日数据(gm-api)
    
    参数：
        - edate: 结束日期，
        - sdate: 开始日期，
        - inclusive: 日期区间开和闭,['both','left','right']
        - timeout: 连接超时时间(默认3s)
        - data_timeout: 连接传输超时时间(默认timeout) 
    '''
    if data_timeout is None: data_timeout = timeout
    now = datetime.datetime.now()
    today = str(now.date())

    if port==443: url = f'https://{ip}/get_daily_valuation'
    else: url = f'http://{ip}:{port}/get_daily_valuation'
    
    if edate is None: edate = today
    
    params = {
        'symbol': symbol,
        'fields': indicators,
        }
    if edate is not None: 
        if inclusive=='left': 
            edate = str(pd.to_datetime(edate)-pd.Timedelta(days=1))
            if sdate is not None and edate<sdate: edate = sdate 
        params = dict(edate=edate,**params)
    if sdate is not None: 
        if inclusive=='right': 
            sdate = str(pd.to_datetime(sdate)+pd.Timedelta(days=1))
            if edate is not None and sdate>edate: sdate = edate 
        if sdate>today: sdate = today 
        if edate is not None and sdate>edate: sdate = edate 
        params = dict(sdate=sdate,**params)
    
    try:
        if retry: json_data = get_url_with_retry(url,params,timeout,data_timeout=data_timeout) 
        else: json_data = get_url_without_retry(url,params,timeout,data_timeout=data_timeout) 
        
        df = pd.DataFrame(json_data)
        if df.shape[0]>0: 
            df.set_index(['symbol','trade_date'],inplace=True)
        else : 
            if debug: print(f" !!! Warning: blank DataFrame for {url}")
        return df
    except (requests.exceptions.RequestException, ValueError, KeyError) as e:
        if debug: print(f" !!! Exception: {e}")
        return pd.DataFrame()  # 返回空的DataFrame

def fetch_gm__kbars_history(symbols,sdate=None,edate=None,tag='1d',inclusive='both',
                    ip='localhost',port=5000,timeout=10,data_timeout=None,retry=True,debug=False):
    '''
    请求历史K线行情数据(gm-api)
    
    参数：
        edate: 结束日期，当给定count时sday参数无效
        sdate: 开始日期，sday=None时为当天日期
        tag: 数据类型，['1m','1d','5m']，分别表示分时、日频、5分钟行情
        inclusive: 日期区间开和闭,['both','left','right']
        timeout: 连接超时时间(默认3s)
        data_timeout: 连接传输超时时间(默认timeout*3)
    '''
    if data_timeout is None: data_timeout = timeout*3

    if port==443: url = f'https://{ip}/get_his'
    else: url = f'http://{ip}:{port}/get_his'

    if isinstance(symbols,list): symbols = ','.join(symbols)
    params = {
        'symbols': symbols,
        'tag': tag,
        }
    
    today = str(pd.Timestamp.today().date())
    if edate is not None: 
        if inclusive=='left': 
            edate = str(pd.to_datetime(edate)-pd.Timedelta(days=1))
            if sdate is not None and edate<sdate: edate = sdate 
        params = dict(edate=edate,**params)
    if sdate is not None: 
        if inclusive=='right': 
            sdate = str(pd.to_datetime(sdate)+pd.Timedelta(days=1))
            if edate is not None and sdate>edate: sdate = edate 
        if sdate>today: sdate = today 
        if edate is not None and sdate>edate: sdate = edate 
        params = dict(sdate=sdate,**params)

    try:
        if retry: json_data = get_url_with_retry(url,params,timeout,data_timeout=data_timeout) 
        else: json_data = get_url_without_retry(url,params,timeout,data_timeout=data_timeout) 
        
        df = pd.DataFrame(json_data)
        if df.shape[0]>0: 
            # df.eob = df.eob.apply(to_localize_asia8)
            df.eob = df.eob.apply(to_drop_utc)
            df.rename(columns={'eob':'timestamp'},inplace=True)
            
            df = df[['symbol','open','high','low','close','volume','timestamp']]
            df.sort_values(['symbol','timestamp'],inplace=True)
            df.set_index(['symbol','timestamp'],inplace=True)
        else : 
            if debug: print(f" !!! Warning: blank DataFrame for {url}")
        return df
    except (requests.exceptions.RequestException, ValueError, KeyError) as e:
        print(f" !!! Exception: {e}")
        return pd.DataFrame()  # 返回空的DataFrame

def fetch_gm__symbol_his_multi_dates(symbol,trade_dates_list,tag='1d',max_connect=50,is_df=False,
                    inclusive='both',ip='localhost',port=5000,timeout=3,data_timeout=30,timesleep=1,debug=False):
    ''' 获取单支股票行情数据，按交易日获取数据, 一次最多50个任务, 主要用于分时行情获取 ''' 
    dfs = {} 
    
    @multitasking.task
    def start(day: str):
        _df = fetch_gm__kbars_history(symbols=symbol,tag=tag,sdate=day,edate=day,
                inclusive=inclusive,timeout=timeout,data_timeout=data_timeout,ip=ip,port=port,debug=debug,)
        if _df.shape[0]>0: 
            dfs[day] = _df.xs(symbol,level='symbol')
    
    if len(trade_dates_list)==0: 
        print(f' !!! Warning: Blank trade_dates_list for fetch_gm__symbol_his_multi_dates() ')
        return pd.DataFrame()
    else:
        sdate = trade_dates_list[0]
        edate = trade_dates_list[-1]
    
    num = 0 
    ncount = 0 
    if debug: print(f'  -> Fetch {symbol}--{tag}-: {sdate} ~ {edate} ... ')
    if debug: print(f'   {num+1:2d}.Fetch ( {ncount:4d} -> task count = {max_connect} ...')
    
    for day in trade_dates_list: 
        ncount += 1 
        ntask = len(multitasking.get_active_tasks())
        if ntask >= max_connect: 
            time.sleep(timesleep) 
            num += 1 
            if debug: print(f'   {num+1:2d}.Fetch ( {symbol}--{tag}- {ncount:4d} -> {day}): task count = {ntask} ...')
        start(day)

    multitasking.wait_for_tasks()

    if is_df:
        num = len(dfs)
        if num>1: 
            df = pd.concat(dfs, axis=0,names=['date'], ignore_index=False) 
            df = df.reset_index(level='date', drop=True) # 去掉date索引
            df = df.sort_index() 
        elif num==1: 
            keys = list(dfs.keys())
            df = dfs[keys[0]] 
        else: df = pd.DataFrame() 
        return df     
        # filted_dfs = {k: v for k, v in dfs.items() if isinstance(v,pd.DataFrame) and v.shape[0]>0 }
        # return pd.concat(filted_dfs, axis=0,names=['date'], ignore_index=False) 
    else:
        return dfs
    
def fetch_kbars__csv_gm_symbol(symbol=None,tag='1m',sdate=None,edate=None,is_index=False, 
        csv_ip='localhost',csv_port=5002,gm_ip='localhost',gm_port=5000,timeout=3,data_timeout=30,debug=False):
    ''' 获取个股行情数据: 先获取csv最新行情, 然后用gm-api接口补全剩余数据 '''
    if symbol is None: 
        print(f' !!! Warning: symbol is a required input for fetch_kbars__csv_gm_symbol(...)')
        return pd.DataFrame()
    
    now = pd.Timestamp.today()
    today = str(now.date())
    if sdate is None: sdate = today
    if edate is None: edate = today
    
    if debug: print(f'  => Get csv.xz for {symbol}--{tag}-: {sdate} ~ {edate} , {csv_ip}:{csv_port}')
    dfdf = fetch_csv__kbars_dates(symbol=symbol,sdate=sdate,edate=edate,tag=tag,
                port=csv_port,ip=csv_ip,timeout=timeout,data_timeout=data_timeout,debug=debug)
    if dfdf.shape[0]==0: 
        dnext = sdate
    else:
        if not isinstance(dfdf.index,pd.Timestamp): 
            dfdf.index = pd.to_datetime(dfdf.index) 
        dnext = str(dfdf.index.date[-1]) 
    
    trade_list = get_trading_days(start_date=dnext,end_date=edate) 
    
    dnext = trade_list[0]
    edate = trade_list[-1]
    dfd1 = pd.DataFrame() 
    if dnext<=edate: 
        if debug: print(f'  -> Get kbars(gm-api) for -{symbol}--{tag}-: {dnext} ~ {edate}, {gm_ip}:{gm_port}')      
        dfd1 = fetch_gm__symbol_his_multi_dates(symbol,trade_list,tag=tag,is_df=True,
                ip=gm_ip,port=gm_port,timeout=timeout,data_timeout=data_timeout,debug=debug)
        if dfd1.shape[0]>0: 
            if is_index and 'volume' in dfd1.columns: 
                dfd1.volume = dfd1.volume/100 
    
    if dfd1.shape[0]>0 and dfdf.shape[0]>0: 
        dfdf = pd.concat([dfdf,dfd1],axis=0) 
        dfdf.drop_duplicates(keep='last',inplace=True)
    elif dfdf.shape[0]==0 and dfd1.shape[0]>0: 
        dfdf = dfd1 
        
    return dfdf


def fetch_csv__kbars_year(symbol=None,year=None,tag='1m',
        ip='localhost',port=5002,timeout=3,data_timeout=None,debug=False): 
    ''' 读取个股某年行情数据(csv.xz文件)
    '''
    if data_timeout is None: data_timeout = timeout*3
    if symbol is None: 
        print(f' sid is a required input parameter.')
        return pd.DataFrame()
    if year is None: year = pd.Timestamp.today().year 
    
    if port==443 : url = f'https://{ip}/download/'
    else: url = f'http://{ip}:{port}/download/'
    
    key    = f'{symbol[:2]}-{symbol[5:7]}'
    # subfld = f'kbars-year/year-{year}--{key}/'
    subfld = f'kbars-year/year-{year}/year-{year}--{key}/'    
    fname  = f'kbars-{tag}--{symbol}--{year}-.csv.xz'
    url = f'{url}{subfld}{fname}'
    try :
        response = requests.get(url, timeout=(timeout,data_timeout))  
        response.raise_for_status()  # 检查是否有 HTTP 错误
        
        dfdf = pd.read_csv(BytesIO(response.content), compression='xz')
        if 'timestamp' in dfdf.columns: 
            dfdf.set_index('timestamp',inplace=True) 
            dfdf.index = pd.to_datetime(dfdf.index) 
        return dfdf 
    except requests.exceptions.Timeout:
        print(f"   Timeout for: {url}")
        return pd.DataFrame()
    except requests.exceptions.RequestException as e:
        if debug: print(f' >>> Exception: {e}')
        return pd.DataFrame()
    
def fetch_csv__kbars_month(symbol=None,year=None,month=None,tag='1m',
        port=5002,ip='localhost',timeout=3,data_timeout=None,debug=False): 
    ''' 读取个股某月行情数据(csv.xz文件)
    '''
    if data_timeout is None: data_timeout = timeout*3
    if symbol is None: 
        print(f' sid is a required input parameter.')
        return pd.DataFrame()
    
    if year is None: year = pd.Timestamp.today().year 
    if month is None: month = pd.Timestamp.today().month 
    
    if port==443 : url = f'https://{ip}/download/'
    else: url = f'http://{ip}:{port}/download/'
    
    key    = f'{symbol[:2]}-{symbol[5:7]}'
    subfld = f'kbars-month/month-{year}/month-{year}-{month:02d}--{key}/'    
    fname  = f'kbars-{tag}--{symbol}--{year}-{month:02d}-.csv.xz'
    url = f'{url}{subfld}{fname}'
    try :
        response = requests.get(url, timeout=(timeout,data_timeout))  
        response.raise_for_status()  # 检查是否有 HTTP 错误
        
        dfdf = pd.read_csv(BytesIO(response.content), compression='xz')
        if 'timestamp' in dfdf.columns: 
            dfdf.set_index('timestamp',inplace=True)
            dfdf.index = pd.to_datetime(dfdf.index)
        return dfdf 
    except requests.exceptions.Timeout:
        print(f"   Timeout for: {url}")
        return pd.DataFrame()
    except requests.exceptions.RequestException as e:
        if debug: print(f' >>> Exception: {e}')
        return pd.DataFrame()


def fetch_csv__kbars_years(symbol=None,syear=None,eyear=None,tag='1m',
        ip='localhost',port=5002,timeout=3,data_timeout=None,debug=False): 
    ''' 读取个股n年行情数据(csv.xz文件) '''
    
    if symbol is None: 
        print(f' sid is a required input parameter.')
        return pd.DataFrame()
    if syear is None: syear = pd.Timestamp.today().year 
    if eyear is None: eyear = pd.Timestamp.today().year 
    
    df = pd.DataFrame()
    for year in range(syear,eyear+1): 
        df1m = fetch_csv__kbars_year(symbol,year=year,tag=tag,
            ip=ip,port=port,timeout=timeout,data_timeout=data_timeout,debug=debug)
        if df1m.shape[0]>0: df = pd.concat([df,df1m],axis=0) if df.shape[0]>0 else df1m 
        elif debug: print(f' >>> Warning: Blank data(csv.xz) for {symbol}--{year}-.')
    return df 

def fetch_csv__kbars_months(symbol=None,smonth=None,emonth=None,tag='1m',
        port=5002,ip='localhost',timeout=3,data_timeout=None,debug=False): 
    ''' 读取个股多月行情数据(csv.xz文件) '''
    if data_timeout is None: data_timeout = timeout*3
    if symbol is None: 
        print(f' sid is a required input parameter.')
        return pd.DataFrame()
    
    year = pd.Timestamp.today().year 
    month = pd.Timestamp.today().month 
    if smonth is None: smonth = f'{year}-{month:02d}'
    if emonth is None: emonth = f'{year}-{month:02d}'
    
    syy = pd.to_datetime(smonth).year
    smm = pd.to_datetime(smonth).month
    eyy = pd.to_datetime(emonth).year
    emm = pd.to_datetime(emonth).month
    
    df = pd.DataFrame() 
    for yy in range(syy,eyy+1): 
        for mm in range(1,13): 
            if yy==syy and mm<smm: continue 
            if yy==eyy and mm>emm: continue 
            
            df1m = fetch_csv__kbars_month(symbol,year=yy,month=mm,tag=tag,
                ip=ip,port=port,timeout=timeout,data_timeout=data_timeout,debug=debug)
            
            if df1m.shape[0]>0: df = pd.concat([df,df1m],axis=0) if df.shape[0]>0 else df1m 
            elif debug: print(f' >>> Warning: Blank data(csv.xz) for {symbol}--{yy}-{mm:02d}.')
            
    return df 


def fetch_csv__kbars_dates(symbol=None,sdate=None,edate=None,tag='1m',
        port=5002,ip='localhost',timeout=3,data_timeout=None,debug=False): 
    ''' 读取个股日期间的行情数据(csv.xz文件): 非开始和结束年读取年份数据，否则读取月份数据 '''
    if data_timeout is None: data_timeout = timeout*3
    if symbol is None: 
        print(f' symbol is a required input parameter.')
        return pd.DataFrame()
    
    today = str(pd.Timestamp.today().date()) 
    if sdate is None: sdate = today 
    if edate is None: edate = today 
    
    syy = pd.to_datetime(sdate).year
    smm = pd.to_datetime(sdate).month
    eyy = pd.to_datetime(edate).year
    emm = pd.to_datetime(edate).month
    
    df = pd.DataFrame() 
    for yy in range(syy,eyy+1): 
        if yy==syy or yy==eyy: 
            smonth = f'{yy}-{smm:02d}' if yy==syy else f'{yy}-01'
            emonth = f'{yy}-{emm:02d}' if yy==eyy else f'{yy}-12'
            df1m = fetch_csv__kbars_months(symbol,smonth=smonth,emonth=emonth,
                ip=ip,port=port,timeout=timeout,data_timeout=data_timeout,debug=debug)    
            
            if df1m.shape[0]>0: df = pd.concat([df,df1m],axis=0) if df.shape[0]>0 else df1m 
            elif debug: print(f' >>> Warning: Blank data(csv.xz) for {smonth}~{emonth}')
        else: 
            df1m = fetch_csv__kbars_year(symbol,year=yy,tag=tag,
                ip=ip,port=port,timeout=timeout,data_timeout=data_timeout,debug=debug)
            
            if df1m.shape[0]>0: df = pd.concat([df,df1m],axis=0) if df.shape[0]>0 else df1m 
            elif debug: print(f' >>> Warning: Blank data(csv.xz) for {symbol}--{yy}-')
    
    if df.shape[0]>0: 
        msk = (df.index>=sdate) & (df.index<=edate+' 15:00')
        df = df.loc[msk]
    return df 

def to_ths_vol(dfvv,df1m,symbol,st_market,fld='./',debug=False):    
    '''导出成交量, 主要用于大盘指数的日成交量和分时成交量'''
    
    row = st_market.loc[symbol]
    sec_name = row['sec_name']
    sec_abbr = row['sec_abbr']
    
    fname = sec_abbr[-5:]+'v.py'
    filename = os.path.join(fld,fname)
    # if debug: print(f'  -> Output path for {symbol}-{sec_name}: {filename}')
    
    sres  = "\n# "+symbol+'-'+sec_name
    sres += "\n# Common used params"
    
    with open(filename, 'w', encoding='utf-8') as file:
        
        file.write(sres)
        file.write("\n#  Data index: [open,close,vol]")
        file.write("\nbarsdata={")
        
        ss24 = '\n}\n\n'
        sres = ''        
        dfvv = dfvv.sort_index(ascending=False)
        for idx, row in  dfvv.iterrows():
            sres = sres+"\n\"%s\" : [%s,%s,%s]," % (
                idx.strftime("%Y%m%d"),
                int(row['open'] if isinstance(row['open'],(int,float)) else 0),
                int(row['close'] if isinstance(row['close'],(int,float)) else 0),
                int(row['volume'] if isinstance(row['volume'],(int,float)) else 0),
                )
            
        file.write(sres)
        file.write(ss24)
        
        file.write("\nbars1m={")        
        ss24 = '\n}\n\n'
        sres = ''        
        df1m = df1m.sort_index(ascending=False)
        for idx, row in  df1m.iterrows():
            sres = sres+"\n\"%s\" : [%s,%s,%s]," % (
                idx.strftime("%Y%m%d %H:%M:%S"),
                int(row['open'] if isinstance(row['open'],(int,float)) else 0),
                int(row['close'] if isinstance(row['close'],(int,float)) else 0),
                int(row['volume'] if isinstance(row['volume'],(int,float)) else 0),
                )
            
        file.write(sres)
        file.write(ss24)        
        
        file.write(ss_ths_1d)

    print('  -> 导出: ',filename)
    
def to_ths_vv1(dfvv,symbol,st_market,fld='./',debug=False):    
    '''导出首尾量成交量和PETTM'''
    
    row = st_market.loc[symbol]
    sec_name = row['sec_name']
    sec_abbr = row['sec_abbr']
    
    fname = sec_abbr[-5:]+'1.py'
    filename = os.path.join(fld,fname)
    # if debug: print(f'  -> Output path for {symbol}-{sec_name}: {filename}')
    
    sres  = "\n# "+symbol+'-'+sec_name
    sres += "\n# Common used params"
    
    with open(filename, 'w', encoding='utf-8') as file:
        file.write(sres)
        file.write("\n#  Data index: [v931,v932,v150]")
        file.write("\nbarsdata={")
        
        ss24 = '\n}\n\n'
        sres = ''
        
        dfvv = dfvv.sort_index(ascending=False)
        for idx, row in  dfvv.iterrows():
            sres = sres+"\n\"%s\" : [%s,%s,%s]," % (
                idx.strftime("%Y%m%d"),
                # int(row['vol']),
                int(row['v931'] if isinstance(row['v931'],(int,float)) else 0),
                int(row['v932'] if isinstance(row['v932'],(int,float)) else 0),
                int(row['v150'] if isinstance(row['v150'],(int,float)) else 0),
                )
            
        file.write(sres)
        file.write(ss24)

        file.write(ss_ths_1m)

    print('  -> 导出: ',filename)
    
def to_ths_vv5(dfvv,symbol,st_market,fld='./',debug=False):    
    '''导出首尾量成交量和PETTM'''
    
    row = st_market.loc[symbol]
    sec_name = row['sec_name']
    sec_abbr = row['sec_abbr']
    
    fname = sec_abbr[-5:]+'5.py'
    filename = os.path.join(fld,fname)
    # if debug: print(f'  -> Output path for {symbol}-{sec_name}: {filename}')
    
    sres  = "\n# "+symbol+'-'+sec_name
    sres += "\n# Common used params"
    
    with open(filename, 'w', encoding='utf-8') as file:
        file.write(sres)
        file.write("\n#  Data index: [v935,v940,pettm]")
        file.write("\nbarsdata={")
        
        ss24 = '\n}\n\n'
        sres = ''
        
        dfvv = dfvv.sort_index(ascending=False)
        # dfvv = dfvv.fillna(0.0)
        for idx, row in  dfvv.iterrows():
            sres = sres+"\n\"%s\" : [%s,%s,%s]," % (
                idx.strftime("%Y%m%d"),
                # int(row['vol']),
                int(row['v935'] if isinstance(row['v935'],(int,float)) else 0),
                int(row['v940'] if isinstance(row['v940'],(int,float)) else 0),
                # int(row['v150']),
                round(row['pe'],2)  if isinstance(row['pe'],(int,float)) else 0.0,
                )
            
        file.write(sres)
        file.write(ss24)

        file.write(ss_ths_5m)

    print('  -> 导出: ',filename)
    
def to_ths_vrd(dfvv,symbol,st_market,fld='./',debug=False):    
    '''导出首尾量比日成交量'''
    
    row = st_market.loc[symbol]
    sec_name = row['sec_name']
    sec_abbr = row['sec_abbr']
    
    fname = sec_abbr[-5:]+'r.py'
    filename = os.path.join(fld,fname)
    # if debug: print(f'  -> Output path for {symbol}-{sec_name}: {filename}')
    
    sres  = "\n# "+symbol+'-'+sec_name
    sres += "\n# Common used params"
    
    with open(filename, 'w', encoding='utf-8') as file:
        file.write(sres)
        file.write("\n#  Data index: [v931,v932,v150]")
        file.write("\nbarsdata={")
        
        ss24 = '\n}\n\n'
        sres = ''
        
        dfvv = dfvv.sort_index(ascending=False)
        for idx, row in  dfvv.iterrows():
            # v9v = row['v931']/row['volume']*100 if row['volume']>0 else 0
            sres = sres+"\n\"%s\" : [%s,%s,%s]," % (
                idx.strftime("%Y%m%d"),
                # round(row['volume'],2),
                round(row['v931'],2),
                round(row['v932'],2),
                round(row['v150'],2),
                # int(v9v if isinstance(v9v,(int,float)) else 0),
                # round(row['pe'],2)  if isinstance(row['pe'],(int,float)) else 0.0,
                )
            
        file.write(sres)
        file.write(ss24)

        file.write(ss_ths_vr)

    print('  -> 导出: ',filename)

def to_ths_vud(dfvv,symbol,st_market,fld='./',debug=False):    
    '''导出首尾量比日成交量'''
    
    row = st_market.loc[symbol]
    sec_name = row['sec_name']
    sec_abbr = row['sec_abbr']
    
    fname = sec_abbr[-5:]+'u.py'
    filename = os.path.join(fld,fname)
    # if debug: print(f'  -> Output path for {symbol}-{sec_name}: {filename}')
    
    sres  = "\n# "+symbol+'-'+sec_name
    sres += "\n# Common used params"
    
    with open(filename, 'w', encoding='utf-8') as file:
        file.write(sres)
        file.write("\n#  Data index: [up,down]")
        file.write("\nbarsdata={")
        
        ss24 = '\n}\n\n'
        sres = ''
        
        dfvv = dfvv.sort_index(ascending=False)
        # dfvv = dfvv.fillna(0.0)
        for idx, row in  dfvv.iterrows():
            up = row['up'] if row['up']>0 else 0
            down = row['down'] if row['down']>0 else 0
            sres = sres+"\n\"%s\" : [%s,%s]," % (
                idx.strftime("%Y%m%d"),
                round(up,2),
                round(down,2),
                )
            
        file.write(sres)
        file.write(ss24)

        file.write(ss_ths_ud)

    print('  -> 导出: ',filename)
    
def to_ths_vpvj(dfvv,symbol,st_market,fld='./',debug=False):    
    '''导出主图指标：量均价'''
    
    row = st_market.loc[symbol]
    sec_name = row['sec_name']
    sec_abbr = row['sec_abbr']
    
    fname = sec_abbr[-5:]+'p.py'
    filename = os.path.join(fld,fname)
    # if debug: print(f'  -> Output path for {symbol}-{sec_name}: {filename}')
    
    sres  = "\n# "+symbol+'-'+sec_name
    sres += "\n# Common used params"
    
    with open(filename, 'w', encoding='utf-8') as file:
        file.write(sres)
        file.write("\n#  Data index: [PVJ,]")
        file.write("\nbarsdata={")
        
        ss24 = '\n}\n\n'
        sres = ''
        
        dfvv = dfvv.sort_index(ascending=False)
        # dfvv = dfvv.fillna(0.0)
        for idx, row in  dfvv.iterrows():
            sres = sres+"\n\"%s\" : [%s]," % (
                idx.strftime("%Y%m%d"),
                round(row['PVJ'],2),
                )
            
        file.write(sres) 
        file.write(ss24) 
        
        file.write(ss_ths_pvj)

    print('  -> 导出: ',filename)
    
def to_ths_cbj(dfvv,symbol,st_market,fld='./',debug=False):    
    '''导出主图指标：成本价'''
    
    row = st_market.loc[symbol]
    sec_name = row['sec_name']
    sec_abbr = row['sec_abbr']
    
    fname = sec_abbr[-5:]+'c.py'
    filename = os.path.join(fld,fname)
    # if debug: print(f'  -> Output path for {symbol}-{sec_name}: {filename}')
    
    sres  = "\n# "+symbol+'-'+sec_name
    sres += "\n# Common used params"
    
    with open(filename, 'w', encoding='utf-8') as file:
        file.write(sres)
        file.write("\n#  Data index: [cbj0,cbj1,cbj2]")
        file.write("\nbarsdata={")
        
        ss24 = '\n}\n\n'
        sres = ''
        
        dfvv = dfvv.sort_index(ascending=False)
        # dfvv = dfvv.fillna(0.0)
        for idx, row in  dfvv.iterrows():
            sres = sres+"\n\"%s\" : [%s,%s,%s]," % (
                idx.strftime("%Y%m%d"),
                round(row['cbj0'],2),
                round(row['cbj1'],2),
                round(row['cbj2'],2),
                )
            
        file.write(sres) 
        file.write(ss24) 
        
        file.write(ss_ths_cbj)

    print('  -> 导出: ',filename)
    
def to_ths_cbf(dfvv,symbol,st_market,fld='./',debug=False):    
    '''导出主图指标：成本价'''
    
    row = st_market.loc[symbol]
    sec_name = row['sec_name']
    sec_abbr = row['sec_abbr']
    
    fname = sec_abbr[-5:]+'f.py'
    filename = os.path.join(fld,fname)
    # if debug: print(f'  -> Output path for {symbol}-{sec_name}: {filename}')
    
    sres  = "\n# "+symbol+'-'+sec_name
    sres += "\n# Common used params"
    
    with open(filename, 'w', encoding='utf-8') as file:
        file.write(sres)
        file.write("\n#  Data index: [cbj0,cbj1,cbj2]")
        file.write("\nbarsdata={")
        
        ss24 = '\n}\n\n'
        sres = ''
        
        dfvv = dfvv.sort_index(ascending=False)
        # dfvv = dfvv.fillna(0.0)
        for idx, row in  dfvv.iterrows():
            sres = sres+"\n\"%s\" : [%s,%s,%s]," % (
                idx.strftime("%Y%m%d"),
                round(row['cbj0'],2),
                round(row['cbj1'],2),
                round(row['cbj2'],2),
                )
            
        file.write(sres) 
        file.write(ss24) 
        
        file.write(ss_ths_cbf)

    print('  -> 导出: ',filename)
    

def get_pvj(df):
    '''由分时数据计算日频-量均价'''
    vol = df['volume']
    hjj = df['hjj'] if 'hjj' in df.columns else df.apply(get_hjj_row,axis=1)
    hvv = hjj*vol
    pvj = hvv.sum() / vol.sum() if vol.sum()>0 else hjj.values[-1]
    ps =  pd.Series(dict(PVJ=pvj))
    return ps

def get_cbj_from_1m(df):
    '''由分时数据计算成本价'''
    hjj = df['HJJ'] if 'HJJ' in df.columns else df.apply(get_hjj_row,axis=1) 
    vol = df['volume']
    if len(vol)<30:
        volmed = vol.mean() # 成交量均值
    elif len(vol)<90:
        volmed = vol.median() # 成交量中值
    else :
        volmed = vol.median()*2 # 成交量中值*2
    
    cbj0 = hjj.iloc[-1]
    cbj1 = cbj0
    cbj2 = cbj0
    
    msk = vol>volmed
    pp = hjj.loc[msk]
    vv = vol.loc[msk]        
    if len(pp)>0:
        cbj0 = (pp*vv).sum()/vv.sum()
        
    msk = pp.index.time<=pd.to_datetime('10:30').time()
    pp1 = pp.loc[msk]
    vv1 = vv.loc[msk]
    if len(pp1)>0 and len(pp1)<len(pp):
        # 早盘到10:30的加权均价
        cbj1 = (pp1*vv1).sum()/vv1.sum()
    
    msk = pp.index.time>pd.to_datetime('10:30').time()
    pp2 = pp.loc[msk]
    vv2 = vv.loc[msk]
    if len(pp2)>0:
        cbj2 = (pp2*vv2).sum()/vv2.sum()
        
    ps =  pd.Series(dict(cbj0=cbj0,cbj1=cbj1,cbj2=cbj2))
    return ps

def get_vol_1m(df):
    '''由分时数据提取日频-成交量'''
    v931 = df.loc[df.index.time == pd.to_datetime('09:31').time()]
    v931 = v931.iloc[0]['volume'] if v931.shape[0]>0 else 0 
    
    v932 = df.loc[df.index.time == pd.to_datetime('09:32').time()]
    v932 = v932.iloc[0]['volume'] if v932.shape[0]>0 else 0 
    
    v150 = df.loc[df.index.time == pd.to_datetime('15:00').time()]
    v150 = v150.iloc[0]['volume'] if v150.shape[0]>0 else 0 
    
    ps =  pd.Series(dict(v931=v931,v932=v932,v150=v150))
    return ps


def get_vol_5m(df):
    '''由分时数据提取日频-成交量'''
    v935 = 0 
    if df.index.time[-1]>pd.to_datetime('09:35').time(): 
        msk = (df.index.time >= pd.to_datetime('09:30').time()) & (df.index.time <= pd.to_datetime('09:35').time())
        v935 = df.loc[msk]
        v935 = v935['volume'].sum() if v935.shape[0]>0 else 0 
        
    v940 = 0 
    if df.index.time[-1]>pd.to_datetime('09:40').time(): 
        msk = (df.index.time >= pd.to_datetime('09:36').time()) & (df.index.time <= pd.to_datetime('09:40').time())
        v940 = df.loc[msk]
        v940 = v940['volume'].sum() if v940.shape[0]>0 else 0 
        
    ps =  pd.Series(dict(v935=v935,v940=v940,))
    return ps

get_kbar = lambda df : pd.Series(dict(
    open   = df['open'].values[0],
    close  = df['close'].values[-1],
    high   = df['high'].max(),
    low    = df['low'].min(),
    volume = df['volume'].sum(),
))
'''由分时数据合成-日K线'''


get_hjj_row = lambda rr: (4*rr['close']+2*rr['open']+rr['high']+rr['low'])/8 
'''由K线计算-黄金价'''


check_open_gg_close_count = lambda df : (df['close'].shift(1)<df['open']).sum()
'''判断开盘价大于收盘价的数量, 单个结果, 可以此数量用于判断交易是否活跃'''

check_open_ll_close_count = lambda df : (df['close'].shift(1)>df['open']).sum()
'''判断开盘价小于收盘价的数量, 单个结果, 可以此数量用于判断交易是否活跃'''

def from_1m(df1m):
    '''从分时行情中提取日频数据'''
    dfgp  = df1m.groupby(pd.to_datetime(df1m.index.date))
    
    df1d  = dfgp.apply(get_kbar,)
    dfpvj = dfgp.apply(get_pvj,)
    df1m  = dfgp.apply(get_vol_1m,)
    df5m  = dfgp.apply(get_vol_5m,)
    dfcb  = dfgp.apply(get_cbj_from_1m,)

    nkbars = 240    
    dfup  = dfgp.apply(check_open_ll_close_count)/nkbars*100
    dfup.name = 'up'
    dfdown = dfgp.apply(check_open_gg_close_count)/nkbars*100
    dfdown.name = 'down'

    dfdf = pd.concat([df1d,dfpvj,dfcb,df1m,df5m,dfup,dfdown],axis=1)
    dfdf['HJJ'] = dfdf.apply(get_hjj_row,axis=1)
    dfdf = dfdf.astype(object).convert_dtypes()
    
    return dfdf 

def read_st_days(datafld='data_ths',ip='localhost',port=5000,debug=False):
    '''读取本地保存的交易日历, 若日历不存在或者最新交易年份小于当前年份则先下载，再保存'''
    # fdd = os.path.join(datafld)
    fdd = datafld 
    
    fname = f'st_calendar.csv.xz'
    fpath = os.path.join(fdd,fname)
    if os.path.exists(fpath): 
        st_calendar = pd.read_csv(fpath,parse_dates=True,index_col=0).apply(pd.to_datetime)
        if st_calendar.shape[0]>0 and st_calendar.index[-1].year>=pd.Timestamp.today().year: 
            return st_calendar 
    
    if debug: print(f'   > Fetch trade calender by gm-api({ip}:{port}) ... ')
    st_calendar = fetch_gm__calendar(ip=ip,port=port,debug=debug)
    st_calendar.to_csv(fpath)
    st_calendar = st_calendar.apply(pd.to_datetime)
        
    return st_calendar

def read_st_market(datafld='data_ths',ip='localhost',port=5000,debug=False):
    '''读取本地保存的市场个股信息, 若文件不存在则先下载，再保存'''
    # fdd = os.path.join(datafld)
    fdd = datafld 
    
    fname = f'st_market.csv.xz'
    fpath = os.path.join(fdd,fname)
    if os.path.exists(fpath): 
        st_market = pd.read_csv(fpath,dtype=str,index_col=0)
        return st_market
    
    if debug: print(f'   > 1. Fetch st_stock by gm-api({ip}:{port}) ... ')
    df_a     = fetch_gm__market_infos(sec='stock',ip=ip,port=port,debug=debug)
    if debug: print(f'   > 2. Fetch st_fund  by gm-api({ip}:{port}) ... ')
    df_etf   = fetch_gm__market_infos(sec='fund',ip=ip,port=port,debug=debug)
    if debug: print(f'   > 3. Fetch st_index by gm-api({ip}:{port}) ... ')
    df_index = fetch_gm__market_infos(sec='index',ip=ip,port=port,debug=debug)
    
    st_market = df_a
    if df_index.shape[0]>0:
        st_market = pd.concat([st_market,df_index],axis=0)
    if df_etf.shape[0]>0:
        st_market = pd.concat([st_market,df_etf],axis=0)
    
    if 'symbol' in st_market.columns: 
        st_market = st_market.set_index('symbol')
        
    if st_market.shape[0]>0:
        st_market.listed_date = st_market.listed_date.dt.date.astype(str)
        st_market.delisted_date = st_market.delisted_date.dt.date.astype(str)
        if debug: print(f'   > Save st_marekt: {fpath}')
        st_market.to_csv(fpath)
    
    return st_market

def read_local_csv_1m(symbol,count=360,fdd='data_ths',ipo=None,save_csv=True,
        gm_ip='localhost',gm_port=5000,csv_ip='localhost',csv_port=5002,debug=False): 
    ''' 读取本地保存的csv分时行情数据 '''
    tag = '1m'
    now = pd.Timestamp.today()
    today = str(now.date())
    edate = now.normalize() 
    sdate = edate-pd.Timedelta(days=count*2)
    trade_list = get_trading_days(str(sdate.date()),str(edate.date()))
    
    trade_list = trade_list[-count:]
    sdate = trade_list[0]
    edate = trade_list[-1]
    
    if ipo is not None: sdate = max(ipo,sdate)
    
    # 确定上一个交易日 
    is_trading_day = today in trade_list 
    pdate = trade_list[-1] 
    if is_trading_day and now<pd.to_datetime('15:00'): 
        pdate = trade_list[-2]
    ptime = pd.to_datetime(pdate+' 15:00')
    
    df = pd.DataFrame()
    fname = f'{symbol}-{tag}.csv.xz'
    fpath = os.path.join(fdd,fname)
    if os.path.exists(fpath):
        print(f'  -> Read local csv({tag}) for {symbol} ...')
        dfoo = pd.read_csv(fpath,parse_dates=True,index_col=0) 
        dfoo.index = pd.to_datetime(dfoo.index) 
        stime_csv = dfoo.index[0]
        etime_csv = dfoo.index[-1]
        # print(f'\tpre_time:{ptime} -> csv-time:{etime_csv}')
        if etime_csv>=ptime: 
            print(f'  >>> No need to update csv.xz({tag}) dataset for {symbol}')
            return dfoo 
        
        df = dfoo 
        ndate = str((etime_csv.normalize()+pd.Timedelta(days=1)).date()) 
        if ndate>today: ndate = today 
    else:
        print(f'  => Local csv is blank for {symbol}, and set the first trading date to fetch -{tag}-(gm-api) ...')
        ndate = sdate 
        
    if is_trading_day and now<pd.to_datetime('15:00'):
        # 未收盘的交易日，将结束日期设置为前一天 
        edate = str((now.normalize()-pd.Timedelta(days=1)).date())
        
    if ndate>edate : ndate = edate 
    
    print(f'  -> To get(gm-api) {symbol}--{tag}-, {ndate} ~ {edate} , {gm_ip}:{gm_port}') 
    dfnn = fetch_kbars__csv_gm_symbol(symbol,tag=tag,sdate=ndate,edate=edate,
                csv_ip=csv_ip,csv_port=csv_port,gm_ip=gm_ip,gm_port=gm_port,debug=debug)
    
    # print(f' \tFetched {tag} data:\n{dfnn}')
    if dfnn.shape[0]>0 and df.shape[0]>0: 
        df = pd.concat([df,dfnn],axis=0)
        df = df.drop_duplicates(keep='last')
    elif df.shape[0]==0 : df = dfnn 
    
    if save_csv and df.shape[0]>0: 
        os.makedirs(fdd,exist_ok=True)
        print(f'  -> Save({tag}): {fpath}')
        df.to_csv(fpath)
        df.index = pd.to_datetime(df.index)
    
    return df 

def read_local_csv_pe(symbol,count=360,fdd='data_ths',ipo=None,save_csv=True,
        gm_ip='localhost',gm_port=5000,csv_ip='localhost',csv_port=5002,debug=False): 
    ''' 读取本地保存的csv分时行情数据 '''
    tag = 'pe'
    now = pd.Timestamp.today()
    today = str(now.date())
    edate = now.normalize() 
    sdate = edate-pd.Timedelta(days=count*2)
    trade_list = get_trading_days(str(sdate.date()),str(edate.date()))
    
    trade_list = trade_list[-count:]
    sdate = trade_list[0]
    edate = trade_list[-1]
    
    if ipo is not None: sdate = max(ipo,sdate)
    
    # 确定上一个交易日 
    is_trading_day = today in trade_list 
    pdate = trade_list[-1] 
    if is_trading_day and now<pd.to_datetime('15:00'): 
        pdate = trade_list[-2]
    # ptime = pd.to_datetime(pdate+' 15:00')
    ptime = pd.to_datetime(pdate)
    
    df = pd.DataFrame()
    fname = f'{symbol}-{tag}.csv.xz'
    fpath = os.path.join(fdd,fname)
    if os.path.exists(fpath):
        print(f'  -> Read local csv({tag}) for {symbol} ...')
        dfoo = pd.read_csv(fpath,parse_dates=True,index_col=0) 
        if dfoo.shape[0]>0: 
            dfoo.index = pd.to_datetime(dfoo.index) 
            stime_csv = dfoo.index[0]
            etime_csv = dfoo.index[-1]
            # print(f'\tpre_time:{ptime} -> csv-time:{etime_csv}')
            if etime_csv>=ptime: 
                print(f'  >>> No need to update csv.xz({tag}) dataset for {symbol}')
                return dfoo 
            
            df = dfoo 
            ndate = str((etime_csv.normalize()+pd.Timedelta(days=1)).date()) 
            if ndate>today: ndate = today 
        else:
            ndate = sdate 
    else:
        print(f'  => Local csv is blank for {symbol}, and set the first trading date to fetch -{tag}-(gm-api) ...')
        ndate = sdate 
        
    print(f'  -> To get(gm-api) {symbol}--{tag}-, {ndate} ~ {today} , {gm_ip}:{gm_port}') 
    # dfnn = fetch_kbars__csv_gm_symbol(symbol,tag=tag,sdate=ndate,edate=today,
    #             csv_ip=csv_ip,csv_port=csv_port,gm_ip=gm_ip,gm_port=gm_port,debug=debug)
    
    dfnn = fetch_gm__daily_valuation(symbol,sdate=ndate,edate=today,ip=gm_ip,port=gm_port,debug=debug)
    if dfnn.shape[0]>0: 
        dfnn = dfnn.xs(symbol,level='symbol')
        dfnn.index = pd.to_datetime(dfnn.index)
    
    # print(f' \tFetched {tag} data:\n{dfnn}')
    if dfnn.shape[0]>0 and df.shape[0]>0: 
        df = pd.concat([df,dfnn],axis=0)
        df = df.drop_duplicates(keep='last')
    elif df.shape[0]==0 : df = dfnn 
    
    if save_csv and df.shape[0]>0: 
        os.makedirs(fdd,exist_ok=True)
        print(f'  -> Save({tag}): {fpath}')
        df.index = df.index.date.astype(str)
        df.to_csv(fpath)
        df.index = pd.to_datetime(df.index)
    
    return df 

def read_local_csv_vv(symbol,count=360,fdd='data_ths',ipo=None,
        gm_ip='localhost',gm_port=5000,csv_ip='localhost',csv_port=5002,debug=False): 
    ''' 读取本地保存的csv日频整理后的行情数据 '''
    tag = 'vv'
    now = pd.Timestamp.today()
    today = str(now.date())
    edate = now.normalize() 
    sdate = edate-pd.Timedelta(days=count*2)
    trade_list = get_trading_days(str(sdate.date()),str(edate.date()))
    
    trade_list = trade_list[-count:]
    sdate = trade_list[0]
    edate = trade_list[-1]
    
    if ipo is not None: sdate = max(ipo,sdate)
    
    # 确定上一个交易日 
    is_trading_day = today in trade_list 
    pdate = trade_list[-1] 
    if is_trading_day and now<pd.to_datetime(today+' 15:00'): 
        pdate = trade_list[-2]
    # ptime = pd.to_datetime(pdate+' 15:00')
    ptime = pd.to_datetime(pdate)
    
    is_blank = False 
    ndate = sdate 
    df = pd.DataFrame()
    fname = f'{symbol}-{tag}.csv.xz'
    fpath = os.path.join(fdd,fname)
    if os.path.exists(fpath):
        print(f' -=> Read local csv({tag}) for {symbol} ...')
        dfoo = pd.read_csv(fpath,parse_dates=True,index_col=0) 
        dfoo.index = pd.to_datetime(dfoo.index) 
        stime_csv = dfoo.index[0]
        etime_csv = dfoo.index[-1]
        print(f'\tpre_time:{ptime} -> csv-time:{etime_csv}')
        if etime_csv>=ptime: 
            print(f'  >>> No need to update csv.xz({tag}) dataset for {symbol}')
            return dfoo 
        
        df = dfoo 
        ndate = str((etime_csv.normalize()+pd.Timedelta(days=1)).date()) 
        if ndate>today: ndate = today 
    else:
        is_blank = True 
        print(f'  => Local csv({tag}) is blank for {symbol}, and set the first trading date to fetch 1m(gm-api) ...')
        df1m = read_local_csv_1m(symbol,count=count,fdd=fdd,ipo=ipo,
                    gm_ip=gm_ip,gm_port=gm_port,csv_ip=csv_ip,csv_port=csv_port,debug=debug)
        if df1m.shape[0]>0: 
            df = from_1m(df1m) 
            ndate = str(df.index[-1].date()) 
    
    if is_trading_day and now<pd.to_datetime('15:00'):
        # 未收盘的交易日，将结束日期设置为前一天 
        edate = str((now.normalize()-pd.Timedelta(days=1)).date())
        
    if ndate>edate : ndate = edate 
    
    print(f'  -> To get(gm-api) {symbol}--1m-: {ndate} ~ {edate} , {gm_ip}:{gm_port}')      
    dfnn = fetch_kbars__csv_gm_symbol(symbol,tag='1m',sdate=ndate,edate=edate,
                csv_ip=csv_ip,csv_port=csv_port,gm_ip=gm_ip,gm_port=gm_port,debug=debug)
    
    if dfnn.shape[0]>0: 
        dfnn = from_1m(dfnn)
    
    if dfnn.shape[0]>0 and df.shape[0]>0: 
        df = pd.concat([df,dfnn],axis=0)
        df = df.drop_duplicates(keep='last')
    elif df.shape[0]==0 : df = dfnn 
    
    if df.shape[0]>0 and ( is_blank or  (not is_trading_day) or (is_trading_day and now>pd.to_datetime('15:00'))): 
        os.makedirs(fdd,exist_ok=True)
        print(f'  -> Save({tag}): {fpath}')
        df.index = df.index.date.astype(str)
        df.to_csv(fpath)
        df.index = pd.to_datetime(df.index)
        
    return df 

def ths_proc_symbols(symbols=None,count=360,tocsv=True,
                    csv_ip='localhost',csv_port=5002,
                    gm_ip='localhost',gm_port=5000,
                    datafld='data_ths',thsfld=None,debug=False):
    ''' 更新股票池内所有个股最新行情 '''
    if thsfld is None: thsfld = 'candle' 
    fdd = datafld
    os.makedirs(fdd,exist_ok=True)
    parent_dir = os.path.dirname(thsfld) 
    
    if symbols is None:
        print(f' >>> Error: symbols is a required input.')
        return False
    
    if not os.path.exists(parent_dir):
        print(f'\n >>> Errors: Parent path for output file directory does not exists: \n\t{parent_dir}\n')
        return 
    
    try: 
        os.makedirs(thsfld,exist_ok=True)
        os.makedirs(os.path.join(thsfld,'Main'),exist_ok=True) 
    except Exception as e:
        print(f'\n >>> Errors: Create output file directory failed: \n {e}\n')
        return 
    
    if symbols is None:
        print(f' >>> Warning: symbols is not provided')
        print(f'   > Set symbols to default: [SHSE.000001,SHSE.601088]')
        st_list = ['SHSE.000001','SHSE.601088'] # ['上证指数','中国神华']
    else: st_list = symbols
    
    #===== 获取交易日历和市场个股信息 =============================================
    # st_days = read_st_days(fld=fld,server_ip=gm_ip,port=gm_port,debug=debug)
    st_market = read_st_market(datafld=datafld,ip=gm_ip,port=gm_port,debug=debug)
    
    df_pool = st_market.loc[st_market.index.isin(st_list)]
    fields = ["sec_name","sec_id","sec_abbr","exchange","listed_date","delisted_date"]
    # if debug: print(f'\n >>> Stock pool: \n {df_pool[fields]}')
    if debug: print(f'\n >>> Stock pool: \n {df_pool}')

    now = pd.Timestamp.today()
    today = str(now.date())
    #======================
    edate = now.normalize()
    sdate = edate-pd.Timedelta(days=count*2)
    trade_list = get_trading_days(str(sdate.date()),str(edate.date()))
    
    trade_list = trade_list[-count:]
    sdate = trade_list[0]
    edate = trade_list[-1]
    
    start_day = trade_list[0]    
    is_trading_day = today in trade_list
    #===== 获取交易日历和市场个股信息 =============================================
    
    num = 0 
    tag = 'vv'
    for sid,row in df_pool.iterrows():
        num += 1 
        print(f'\n -=> {num:2d}. symbol({sid}-{row["sec_name"]}): ipo={row["listed_date"]}')
        ipo = row['listed_date']
        is_stock = int(row['sec_type1'])==1010
        # sday = start_day if start_day>ipo else ipo 

        is_his_local = True 
        
        # 1. 先读取本地保存的数据
        fname = f'{sid}-{tag}.csv.xz'
        fpath = os.path.join(fdd,fname)
        dfoo = read_local_csv_vv(sid,count=count,fdd=datafld,ipo=ipo, 
                gm_ip=gm_ip,gm_port=gm_port,csv_ip=csv_ip,csv_port=csv_port,debug=debug)
        
        #========= 将数据保存到硬盘，下次调用时直接读取历史数据，无需重新下载已经下载过的数据。
        dfvv = dfoo 
        dfvv['pe'] = 0 
        if is_stock: 
            dfpe = read_local_csv_pe(sid,fdd=fdd,gm_ip=gm_ip,gm_port=gm_port,csv_ip=csv_ip,csv_port=csv_port,debug=debug)
            if dfpe.shape[0]>0: dfvv['pe'] = dfpe ['pe_ttm']
            # print(f' dfpe:\n{dfpe}')
            
        is_closed = pd.Timestamp.today()>pd.to_datetime('15:00')
        if is_trading_day and not is_closed: 
            # 当日为交易日，同时市场未收盘，则获取最新分时行情 
            print(f'  -> To get latest kbars(1m) for ({sid}-{row["sec_name"]}--{tag}-), >> {today} << ({gm_ip}:{gm_port})')      
            
            _df = fetch_gm__kbars_history(symbols=sid,tag='1m',sdate=today,edate=today,
                    ip=gm_ip,port=gm_port,debug=debug,)
            if _df.shape[0]>0: 
                _df = _df.xs(sid,level='symbol')
                _vv = from_1m(_df)
                
                pe_prev = dfvv['pe'].values[-1]
                cc_prev = dfvv['close'].values[-1]
                _pe = _vv['close']*pe_prev/cc_prev if cc_prev>0 else 0
                _vv['pe'] = _pe if _pe is None else 0
                # print(f' _vv:\n{_vv}')
                
                idx = _vv.index.values[-1]
                dfvv.loc[idx] = _vv.loc[idx]
        
        # dfvv = dfvv.convert_dtypes()
        # print(f' \t=> dfvv:\n{dfvv}')
        # break # Test 
        
        if dfvv.shape[0]>0: 
            # os.makedirs(fdd,exist_ok=True)
            # 如果是交易日，只有在收盘之后才保存数据 
            dfvv = dfvv.fillna(0)
            print(f'  -> 保存行情数据: {fpath}')
            if tocsv: 
                dfvv.index.name = 'date'
                # dfvv.to_csv(fpath)
                if is_trading_day and is_closed:
                    dfvv.to_csv(fpath)
                elif not is_his_local:
                    dfoo.to_csv(fpath)

            #========= 导出同花顺指标 =========================================
            to_ths_vv1(dfvv,symbol=sid,st_market=df_pool,fld=thsfld,debug=debug)
            to_ths_vv5(dfvv,symbol=sid,st_market=df_pool,fld=thsfld,debug=debug)
            to_ths_vrd(dfvv,symbol=sid,st_market=df_pool,fld=thsfld,debug=debug)
            to_ths_vud(dfvv,symbol=sid,st_market=df_pool,fld=thsfld,debug=debug)
            to_ths_cbf(dfvv,symbol=sid,st_market=df_pool,fld=thsfld,debug=debug)
            
            to_ths_vpvj(dfvv,symbol=sid,st_market=df_pool,fld=os.path.join(thsfld,'Main'),debug=debug)
            
            if sid not in INDEX_LIST: 
                to_ths_cbj(dfvv,symbol=sid,st_market=df_pool,fld=os.path.join(thsfld,'Main'),debug=debug)
            
            df1m = None 
            if sid in ['SHSE.000001']: 
                # 对于指数行情，应额外输出一下日成交量和分时成交量，
                #   用于看盘，分析个股与大盘之间的强弱分析 
                # df1m = fetch_gm__symbol_his_multi_dates(sid,trade_list[-10:],tag='1m',
                #             ip=gm_ip,port=gm_port,debug=debug)
                df1m = read_local_csv_1m(sid,count=count,fdd=fdd,
                    gm_ip=gm_ip,gm_port=gm_port,csv_ip=csv_ip,csv_port=csv_port,debug=debug)
                
                df1m = df1m.loc[df1m.index>=trade_list[-10]]
                # print(df1m) 
                to_ths_vol(dfvv,df1m,symbol=sid,st_market=df_pool,fld=thsfld,debug=debug) 
            
        # break # Test 
    
#==========================
if __name__ == "__main__":
    import sys,datetime 
    
    debug = True    
    args = sys.argv[1:] # sys.argv[0] 是脚本名，后面是参数
    
    #=======================================
    gm_ip    = 'gm.zhwdk.com'   #
    gm_port  = 443
    
    csv_ip   = 'localhost'
    csv_ip   = gm_ip
    csv_port = 5002
    
    ths_dir  = None 
    data_fld = 'candle'
    
    count = 360 
    nyear = 2 
    symbols = [ 
        # 'SHSE.000001', # 上证指数
        'SHSE.601088', # 中国神华
        ]
    
    with open("cfg.toml", "r",encoding='utf-8') as f:
        cfg = toml.load(f)        
        print(f' -=> Loading params from: cfg.toml')
        
        debug     = cfg['ths']['debug']
        ths_dir   = cfg['ths']['ths_dir']
        data_fld  = cfg['ths']['data_dir']
        
        gm_ip    = cfg['ths']['gm_ip']
        gm_port  = int(cfg['ths']['gm_port'])
        csv_ip   = cfg['ths']['csv_ip']
        csv_port = int(cfg['ths']['csv_port'])
        
        count      = int(cfg['ths']['count'])
        symbols    = cfg['ths']['symbols']
        INDEX_LIST = cfg['ths']['index_list']
    
    # 创建THS指标文件夹 
    if ths_dir and not os.path.exists(ths_dir): os.makedirs(ths_dir,exist_ok=True) 
    
    # 只有交易日才进行数据的更新
    #==================================
    if len(args)>0 and args[0]=='test': 
        print(f'\n >>> Start to test script ...')         
        year = 2024 
        
        symbol = 'SHSE.601088'
        # symbol = 'SHSE.000001'
        ipo = '2007-10-09'
        # ipo = '2005-01-01'
        
        df1m = read_local_csv_1m(symbol,count=count,fdd=data_fld,ipo=ipo,save_csv=True,
                    gm_ip=gm_ip,gm_port=gm_port,csv_ip=csv_ip,csv_port=csv_port,debug=debug)
        print(f'df1m:\n{df1m}')

        dfpe = read_local_csv_pe(symbol,count=count,fdd=data_fld,ipo=ipo,
                    gm_ip=gm_ip,gm_port=gm_port,csv_ip=csv_ip,csv_port=csv_port,debug=debug)
        print(f'dfpe:\n{dfpe}')
        
        dfvv = read_local_csv_vv(symbol,count=count,fdd=data_fld,ipo=ipo,
                    gm_ip=gm_ip,gm_port=gm_port,csv_ip=csv_ip,csv_port=csv_port,debug=debug)
        print(f'dfvv:\n{dfvv}')
        
        
    elif  len(args)>0 and args[0]=='update': 
        print(f'\n >>> Start update THS daily indications ...')        
        ths_proc_symbols(symbols=symbols,count=count,tocsv=True,
                        csv_ip=csv_ip,csv_port=csv_port,gm_ip=gm_ip,gm_port=gm_port,
                        datafld=data_fld,thsfld=ths_dir,debug=debug) 
    else:
        print(f'\n >>> Opps: cmd options or params error: {args}  ...')
        print(f'\n         > python ths_v15.py test ')
        print(f'\n         > python ths_v15.py update ')
        
        
    #============================
    print(f'\n !!! Nice, All misson finished at {datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}. \n')
    
    