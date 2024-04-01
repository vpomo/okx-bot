https://tproger.ru/translations/deploy-a-secure-golang-rest-api

1. POST - http://localhost:8000/api/user/new
   {
   "email": "pomo@gmail.com",
   "password": "rrr444@@"
   }
2. POST - http://localhost:8000/api/user/login
   {
   "email": "pomo@gmail.com",
   "password": "rrr444@@"
   }
3. POST - http://localhost:8000/api/contacts/new
   {
   "name": "Иван",
   "phone": "981-082-099"
   }
4. GET - http://localhost:8000/api/me/contacts


/build/rest