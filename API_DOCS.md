# Trip History API Documentation

เอกสารแนะนำการใช้งาน API สำหรับระบบบันทึกประวัติการเดินทางและการติดตามพิกัด (Trip History & Tracking)

---

## Base URL

\`\`\`text
http://localhost:8080
\`\`\`

---

## 1. การจัดการผู้ใช้งาน (User Management)

### 1.1 สมัครสมาชิก (Register User)

สร้างบัญชีผู้ใช้งานใหม่โดยบันทึกลงใน Firebase Authentication (จะ Hash รหัสผ่านให้อัตโนมัติ) และบันทึก Profile ลงใน Firestore

- **URL**: `/api/users/register`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
    "email": "driver1@example.com",
    "password": "securepassword123",
    "name": "Somchai Jaidee",
    "role": "driver"
  }
  ```
  _(หมายเหตุ: `role` สามารถส่งเป็น `admin` หรือ `driver` ก็ได้ หากไม่ส่งระบบจะให้ค่าเริ่มต้นเป็น `driver`)_
- **Success Response** (200 OK):
  ```json
  {
    "status": "User registered",
    "userId": "wXyZ1234abc...",
    "verificationLink": "https://identitytoolkit.googleapis.com/v1/action?mode=verifyEmail&oobCode=..."
  }
  ```

### 1.2 ลบบัญชีผู้ใช้งาน (Delete User)

ลบผู้ใช้ออกจากระบบ (ทั้งใน Firestore และ Firebase Auth)

- **URL**: `/api/users/:id`
  - _:id_ คือ `userId` ของผู้ใช้ที่ต้องการลบ
- **Method**: `DELETE`
- **Success Response** (200 OK):
  ```json
  {
    "status": "User deleted"
  }
  ```

### 1.3 ดึงข้อมูลผู้ใช้ทั้งหมด (Get All Users)

ดึงรายชื่อผู้ใช้ทั้งหมดในระบบ (สำหรับฝั่ง Admin)

- **URL**: `/api/users`
- **Method**: `GET`
- **Success Response** (200 OK):
  ```json
  [
    {
      "userId": "wXyZ1234abc...",
      "email": "driver1@example.com",
      "name": "Somchai Jaidee",
      "role": "driver"
    },
    {
      "userId": "aBcD5678xyz...",
      "email": "admin@example.com",
      "name": "Admin System",
      "role": "admin"
    }
  ]
  ```

---

## 2. การจัดการทริป (Trip Management)

### 2.1 เริ่มการเดินทาง (Start Trip)

เปิดทริปใหม่ สถานะจะเป็น `ongoing`

- **URL**: `/api/trips/start`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
    "userId": "wXyZ1234abc...",
    "planId": "plan456_optional"
  }
  ```
- **Success Response** (200 OK):
  ```json
  {
    "message": "Trip started successfully",
    "tripId": "tripId_001"
  }
  ```

### 2.2 จบการเดินทาง (End Trip)

บันทึกเวลาจบ ระยะทางรวม และรูปภาพยืนยันการจบงาน

- **URL**: `/api/trips/:id/end`
  - _:id_ คือ `tripId` ของทริปนั้น
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
    "distance": 15.5,
    "imageUrl": "https://example.com/images/end_trip.jpg"
  }
  ```
- **Success Response** (200 OK):
  ```json
  {
    "status": "Trip ended"
  }
  ```

### 2.3 ดึงประวัติการเดินทาง (Get Trips)

ดึงรายการทริปทั้งหมดของ User คนนั้น

- **URL**: `/api/trips?userId={userId}`
- **Method**: `GET`
- **Success Response** (200 OK):
  ```json
  [
    {
      "userId": "wXyZ1234abc...",
      "status": "completed",
      "startTime": "2026-09-03T10:00:00Z",
      "endTime": "2026-09-03T11:30:00Z",
      "distance": 15.5,
      "imageUrl": "https://example.com/images/end_trip.jpg"
    }
  ]
  ```

### 2.4 ดึงประวัติการเดินทางตามช่วงเวลา (Get Trip History by Time Range)

ดึงรายการทริปทั้งหมดของ User คนนั้น ที่เริ่มต้นในช่วงเวลาที่กำหนด (ต้องสร้าง Composite Index บน Firestore ก่อนใช้งาน)

- **URL**: `/api/history?userId={userId}&start={startTime}&end={endTime}`
- **Method**: `GET`
- **Query Params**:
  - `userId` (required): รหัสผู้ใช้งาน
  - `start` (required): เวลาเริ่มต้น รูปแบบ RFC3339 เช่น `2026-09-03T00:00:00Z`
  - `end` (required): เวลาสิ้นสุด รูปแบบ RFC3339 เช่น `2026-09-03T23:59:59Z`
- **Success Response** (200 OK):
  ```json
  [
    {
      "id": "abc123tripId...",
      "planId": "plan456...",
      "userId": "wXyZ1234abc...",
      "status": "completed",
      "startTime": "2026-09-03T10:00:00Z",
      "endTime": "2026-09-03T11:30:00Z",
      "distance": 15.5,
      "imageUrl": "https://example.com/images/end_trip.jpg"
    }
  ]
  ```

---

## 3. การติดตามและเส้นทาง (Tracking & Routes)

### 3.0 แจ้งปัญหาในการเดินทาง (Report Issue)

คนขับสามารถส่งข้อมูลรายงานปัญหาที่พบเจอระหว่างเดินทาง เพื่อผูกติดกับทริปได้

- **URL**: `/api/trips/:id/issues`
  - _:id_ คือ `tripId`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
    "title": "รถยางแตก",
    "description": "เกิดอุบัติเหตุยางรั่วที่ถนนบางนา-ตราด"
  }
  ```
- **Success Response** (200 OK):
  ```json
  {
    "status": "Issue reported successfully"
  }
  ```

### 3.1 ส่งข้อมูลพิกัด (Update Location)

ส่งพิกัดระหว่างเดินทาง ข้อมูลนี้จะถูกบันทึกลงประวัติ (Firestore) และ อัปเดตไปยัง Redis พร้อมทั้ง **ยิง WebSocket บรอดแคสต์แบบเรียลไทม์** ทันที

- **URL**: `/api/trips/:id/locations`
  - _:id_ คือ `tripId`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
    "userId": "wXyZ1234abc...",
    "lat": 13.7563,
    "lng": 100.5018,
    "speed": 60.5
  }
  ```
- **Success Response** (200 OK):
  ```json
  {
    "status": "Location updated"
  }
  ```

### 3.2 ดึงพิกัดเส้นทางเพื่อวาดแผนที่ (Get Trip Route)

ดึงประวัติพิกัดทั้งหมดของทริป เพื่อนำไปวาดเส้น Polyline บนแผนที่

- **URL**: `/api/trips/:id/route`
- **Method**: `GET`
- **Success Response** (200 OK):
  ```json
  [
    {
      "lat": 13.7563,
      "lng": 100.5018,
      "speed": 60.5,
      "timestamp": "2026-09-03T10:15:00Z"
    },
    ...
  ]
  ```

### 3.3 การติดตามเรียลไทม์ผ่าน WebSocket (Real-time Web Tracking)

ใช้สำหรับฝั่งหน้าเว็บ (เช่น Next.js) ต่อเข้ามาฟังการเคลื่อนไหวของรถแบบเรียลไทม์

- **URL**: `/ws/tracking/:userId`
  - หมายเหตุ: หากหน้าเว็บแอดมินต้องการดูข้อมูลของ **รถทุกคันพร้อมกัน** ให้ใช้ `:userId` เป็นคำว่า `all` เช่น `/ws/tracking/all`
  - _:userId_ คือ `userId` ที่กำลังเดินทาง หรือพิมพ์คำว่า `all` เพื่อดึงพิกัดของคนขับทุกคน (สำหรับ Admin)
- **การเชื่อมต่อ**: ใช้ WebSocket Client (เช่น ใน Frontend ใช้ `new WebSocket('ws://localhost:8080/ws/tracking/all')`)
- **การรับข้อมูล**: เมื่อ Server มีการอัปเดต จะส่งข้อความกลับมาเป็น JSON:
  ```json
  {
    "userId": "wXyZ1234abc...",
    "lat": 13.7563,
    "lng": 100.5018,
    "speed": 60.5,
    "timestamp": "0001-01-01T00:00:00Z"
  }
  ```

---

## 4. เหตุการณ์และข้อผิดพลาด (Events & Logs)

### 4.1 เช็กอินตามจุด (Check-in)

บันทึกเหตุการณ์ต่างๆ ระหว่างทริป เช่น ถึงจุดนัดหมาย, ส่งของเรียบร้อย

- **URL**: `/api/trips/:id/checkin`
- **Method**: `POST`
- **Body** (JSON):
  \`\`\`json
  {
  "lat": 13.7600,
  "lng": 100.5100,
  "timestamp": "2026-09-03T10:30:00Z",
  "message": "Arrived at customer location",
  "imageUrl": "https://example.com/images/checkin.jpg"
  }
  \`\`\`
- **Success Response** (200 OK):
  \`\`\`json
  {
  "status": "Check-in saved"
  }
  \`\`\`

### 4.2 บันทึก Log ระบบ (Save System Log)

ส่งข้อมูลปัญหาของระบบ เช่น GPS สัญญาณหาย, อินเทอร์เน็ตหลุด (สำหรับวิเคราะห์ปัญหาภายหลัง)

- **URL**: `/api/logs`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
    "userId": "wXyZ1234abc...",
    "eventType": "GPS_LOST",
    "message": "GPS signal lost in tunnel",
    "timestamp": "2026-09-03T10:45:00Z"
  }
  ```
- **Success Response** (200 OK):
  ```json
  {
    "status": "Log saved"
  }
  ```

---

## 5. จัดการแผนงาน (Job Plans)

### 5.1 สร้างแผนงานใหม่ (Create Job Plan)

สำหรับ Admin ใช้เพื่อสร้างแผนงาน (Route Plan) และจ่ายงานให้คนขับ

- **URL**: `/api/plans`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
    "driverId": "wXyZ1234abc...",
    "description": "ส่งพัสดุเส้นทาง ลาดพร้าว-รามอินทรา",
    "startTime": "2026-09-04T08:00:00Z",
    "pickupPoint": {
      "lat": 13.7563,
      "lng": 100.5018,
      "address": "คลังสินค้า A"
    },
    "dropOffPoints": [
      {
        "lat": 13.8055,
        "lng": 100.5746,
        "address": "ลูกค้า 1"
      },
      {
        "lat": 13.8211,
        "lng": 100.6543,
        "address": "ลูกค้า 2"
      }
    ]
  }
  ```
- **Success Response** (200 OK):
  ```json
  {
    "message": "Job plan created successfully",
    "planId": "plan_001"
  }
  ```

### 5.2 ดึงแผนงานทั้งหมด (Get Job Plans)

ดึงแผนงานทั้งหมด หรือดึงเฉพาะของคนขับแต่ละคน

- **URL**: `/api/plans` (หรือ `/api/plans?driverId={userId}` เพื่อดูเฉพาะของคนขับคนนั้น)
- **Method**: `GET`
- **Success Response** (200 OK):
  ```json
  [
    {
      "planId": "plan_001",
      "driverId": "wXyZ1234abc...",
      "description": "ส่งพัสดุเส้นทาง ลาดพร้าว-รามอินทรา",
      "status": "pending",
      "createdAt": "2026-09-03T17:00:00Z"
    }
  ]
  ```

### 5.3 อัปเดตสถานะแผนงาน (Update Job Plan Status)

อัปเดตเมื่อคนขับเริ่มงาน (in_progress) หรือจบงาน (completed)

- **URL**: `/api/plans/:id/status`
  - _:id_ คือ `planId`
- **Method**: `PATCH`
- **Body** (JSON):
  ```json
  {
    "status": "in_progress"
  }
  ```
- **Success Response** (200 OK):
  ```json
  {
    "status": "Job plan updated"
  }
  ```

---

## 6. การจัดการระบบ (System Management)

### 6.1 ล้างข้อมูลทั้งหมด (Clear All Data)

ใช้สำหรับล้างข้อมูลในระบบเพื่อรีเซ็ตก่อนเริ่มใช้งานจริง (จะลบ Users ใน Firebase Auth, ลบข้อมูลใน Firestore และลบพิกัดใน Redis ทิ้งทั้งหมด)

- **URL**: `/api/system/clear-data`
- **Method**: `POST`
- **Success Response** (200 OK):
  ```json
  {
    "status": "All data cleared successfully"
  }
  ```

### 6.2 ลบผู้ใช้งาน (Delete User)

- **URL**: `/api/users/:id`
- **Method**: `DELETE`
- **Success Response** (200 OK):
  ```json
  {
    "message": "user deleted successfully"
  }
  ```

### 6.3 อัปเดต FCM Token สำหรับแจ้งเตือน (Update FCM Token)

แอปมือถือส่ง Token ของเครื่องมาให้ Server เพื่อใช้รับ Push Notification (FCM)

- **URL**: `/api/users/:id/fcm-token`
  - _:id_ คือ `userId`
- **Method**: `PATCH`
- **Body** (JSON):
  ```json
  {
    "fcmToken": "eKxT_abc123..."
  }
  ```
- **Success Response** (200 OK):
  ```json
  {
    "status": "token updated"
  }
  ```
