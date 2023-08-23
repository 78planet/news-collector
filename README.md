# news-collector

아침에 잠에서 일어나 핸드폰을 집어들어 네이버를 들어가보니 인터넷에 접속이 안된다. 밤새 무슨 일이 생긴걸까? 집에는 티비와 라디오도 없다. 

인터넷에 연결을 못할 때 최신의 뉴스를 확인해 보자. 

네이버 뉴스 메인페이지의 모든 뉴스들을 HTML 파일로 저장한다. 

저장된 HTML 뉴스 파일들을 파일탐색기에서 확인 가능하다. 

1시간마다 뉴스들을 수집한다. 

```
0 * * * * cd /users/will/newsarchive && ./collector >> /Users/will/newsArchive/test.log
```
