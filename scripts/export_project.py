#coding=utf-8
#导出未完结项目到Excel表格
import MySQLdb
import xlwt
import sys

reload(sys)
sys.setdefaultencoding('utf-8')

conn = MySQLdb.connect(host = "127.0.0.1", user = "root", passwd = "", db = "emsdatabase", port = 3306, charset = "utf8")
cursor = conn.cursor()
#导出只有最后一步未完结的项目
#sqlText = 'SELECT * FROM tbs001_developprojectbasicinfo WHERE (xiangmukaopingchuliren IS NOT NULL || xiangmukaopingchuliriqi IS NOT NULL) AND (reportstate IS NULL || ReportingTrace IS NULL)'
#导出所有未完结项目
sqlText = 'SELECT * FROM tbs001_developprojectbasicinfo WHERE shoulituihuishanchu=1 AND (reportstate IS NULL || ReportingTrace IS NULL)'
cursor.execute(sqlText)
data = cursor.fetchall()
dataCount = cursor.rowcount
print dataCount

sourceFile = 'expro.xls'
workBook = xlwt.Workbook()
sheet = workBook.add_sheet('sheet1', cell_overwrite_ok=True)
sheet.write(0, 0, u'项目名称')
sheet.write(0, 1, u'项目ID')
sheet.write(0, 2, u'建设地点')
sheet.write(0, 3, u'项目联系人')
sheet.write(0, 4, u'项目受理人')
sheet.write(0, 5, u'项目受理时间')
sheet.write(0, 6, u'部门')
sheet.write(0, 7, u'项目描述')
sheet.write(0, 8, u'建设单位ID')
sheet.write(0, 9, u'评价单位ID')
sheet.write(0, 10, u'当前流程')
i = 0
for i in range(0, dataCount):
    print u'项目名称:', data[i][0]
    print u'项目ID：', data[i][1]
    print u'建设地点：', data[i][2]
    print u'人：', data[i][4]
    print u'项目受理人：', data[i][56]
    print u'项目受理时间：', data[i][57]
    print u'部门：', data[i][53]
    print u'项目描述：', data[i][102]
    print u'建设单位ID：', data[i][100]
    print u'评价单位ID：', data[i][101]
    print u'当前流程：', data[i][139]

    sheet.write(i+1, 0, data[i-1][0])
    sheet.write(i+1, 1, data[i-1][1])
    sheet.write(i+1, 2, data[i-1][2])
    sheet.write(i+1, 3, data[i-1][4])
    sheet.write(i+1, 4, data[i-1][56])
    sheet.write(i+1, 5, data[i-1][57])
    sheet.write(i+1, 6, data[i-1][53])
    sheet.write(i+1, 7, data[i-1][102])
    sheet.write(i+1, 8, data[i-1][100])
    sheet.write(i+1, 9, data[i-1][101])
    sheet.write(i+1, 10, data[i-1][139])
workBook.save(sourceFile)

cursor.close()
conn.close()