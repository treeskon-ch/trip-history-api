# Treesukon TMS (Trip History System)

ระบบบริหารจัดการและติดตามสถานะการจัดส่งสินค้า (Transportation Management System) ที่ประกอบด้วยแอปพลิเคชันมือถือสำหรับคนขับ และเว็บแดชบอร์ดสำหรับผู้ดูแลระบบ เพื่อดูประวัติการเดินทางและติดตามตำแหน่งแบบเรียลไทม์

---

## ลิงก์สำหรับทดสอบระบบ (Testing URLs)

- **Web Dashboard (สำหรับ Admin):** `https://trip-history-web.vercel.app` _(หรือ URL ที่ได้จากการ Deploy บน Vercel)_
- **Backend API (Base URL):** `https://trip-history-api.onrender.com`
- **WebSocket URL:** `wss://trip-history-api.onrender.com/ws/tracking/all`

---

## วิธีการติดตั้งและใช้งานแอปพลิเคชัน (Mobile App)

สำหรับคนขับรถ สามารถติดตั้งแอปพลิเคชันผ่านไฟล์ APK ได้ตามขั้นตอนดังนี้:

1.  ดาวน์โหลดไฟล์ `trip-history.apk` ลงในสมาร์ทโฟน Android
2.  ไปที่ **Settings (การตั้งค่า) > Security (ความปลอดภัย)** และเปิดใช้งาน **"Install unknown apps" (อนุญาตให้ติดตั้งแอปที่ไม่รู้จัก)**
3.  เปิดไฟล์ `trip-history.apk` ที่ดาวน์โหลดมา และกด **Install (ติดตั้ง)**
4.  เมื่อติดตั้งเสร็จสิ้น ให้เปิดแอปและเข้าสู่ระบบด้วยบัญชีคนขับรถ (Driver Account)
5.  กด **"เช็คอิน"** เพื่อเริ่มบันทึกการเดินทาง ระบบจะเริ่มส่งพิกัด GPS แบบเรียลไทม์
6.  สามารถกด **"แจ้งปัญหา"** ได้ระหว่างทาง และกด **"จบงาน"** เมื่อถึงที่หมาย

---

## โครงสร้างเทคโนโลยีและสถาปัตยกรรม (Architecture & Technologies)

โปรเจคนี้ออกแบบภายใต้สถาปัตยกรรม **Client-Server & Event-Driven Architecture** เพื่อให้รองรับการแสดงผลพิกัดและแจ้งเตือนแบบ Real-time โดยแยกส่วนประกอบหลักดังนี้:

### 1. Technology & Frameworks

- **Mobile App (Driver):** `Flutter` & `Dart` (ใช้แพ็กเกจ Geolocator สำหรับดึงพิกัด และ HTTP สำหรับยิง API)
- **Web Dashboard (Admin):** `Next.js` (React), `TailwindCSS` (สำหรับ UI), `Leaflet.js` (สำหรับการแสดงแผนที่)
- **Backend API:** `Golang` ร่วมกับ `Gin Framework` โดดเด่นเรื่องประสิทธิภาพสูงและการจัดการ Concurrency

### 2. Database & Caching

- **Primary Database:** `Firebase Firestore (NoSQL)`
  - เหมาะสำหรับการเก็บข้อมูลที่มีลักษณะเป็น Document และ Scale ได้ง่าย (ไม่ต้องทำ Schema Migration)
- **In-Memory Cache:** `Upstash Redis`
  - ใช้สำหรับเก็บสถานะคนขับที่กำลังออนไลน์ และพิกัดล่าสุด เพื่อลดภาระการอ่าน/เขียน (Read/Write) ลง Firestore โดยตรง

### 3. AI Assistance

- **AI Coding Assistant:** พัฒนาและออกแบบโครงสร้างโค้ดร่วมกับ AI Agent เพื่อเพิ่มความเร็วในการขึ้นโครงโปรเจค (Rapid Prototyping) การแก้บั๊ก (Debugging) และสร้าง UI/UX ที่ทันสมัย

---

## การรับและจัดเก็บข้อมูลการเดินทาง (Data Collection & Database Design)

ระบบฐานข้อมูลถูกออกแบบแบบ **NoSQL (Firestore)** เพื่อความยืดหยุ่นในการจัดเก็บข้อมูล โดยมีโครงสร้าง Collection หลักๆ ดังนี้:

### `users` Collection

เก็บข้อมูลผู้ใช้งาน ทั้งคนขับและแอดมิน

- `userId` (String) - รหัสผู้ใช้
- `name` (String) - ชื่อ-นามสกุล
- `role` (String) - `driver` หรือ `admin`

### `trips` Collection

เก็บบันทึกประวัติการเดินทางในแต่ละรอบ (1 ทริป / 1 Document)

- `tripId` (String) - รหัสการเดินทาง
- `userId` (String) - รหัสคนขับที่อ้างอิงจาก `users`
- `status` (String) - `active` (กำลังวิ่ง), `completed` (จบงาน)
- `startTime` (Timestamp) - เวลาที่กดเช็คอิน
- `endTime` (Timestamp) - เวลาที่กดจบงาน
- `distance` (Number) - ระยะทางรวมของทริปนี้ (กิโลเมตร)

### `locations` Collection (Sub-collection ของ `trips`)

เก็บจุดพิกัด GPS ถี่ๆ ตามเวลาที่แอปมือถือส่งมา

- `lat` (Number) - ละติจูด
- `lng` (Number) - ลองจิจูด
- `speed` (Number) - ความเร็ว (กม./ชม.)
- `timestamp` (Timestamp) - เวลาที่บันทึกพิกัด

### `issues` Collection

เก็บรายงานปัญหาที่คนขับส่งเข้ามา

- `tripId` (String) - อ้างอิงทริป
- `title` (String) - หัวข้อปัญหา (เช่น ยางแตก, รถติดหนัก)
- `description` (String) - รายละเอียดเพิ่มเติม
- `reportedAt` (Timestamp) - เวลาที่แจ้งปัญหา

---

## Workflow การทำงานของระบบ

1.  **Mobile App** ยิง API `/api/trips/checkin` เพื่อเริ่มงาน
2.  แอปจะดึงพิกัด GPS มือถือและส่งไปที่ `/api/trips/:id/locations`
3.  **Golang Backend** จะรับพิกัดมา:
    - บันทึกลง Firestore
    - อัปเดตสถานะล่าสุดใน Redis
    - กระจายข้อมูล (Broadcast) ผ่าน **WebSocket** ไปยังผู้ใช้งานทุกคนที่เชื่อมต่ออยู่
4.  **Web Dashboard** ที่เชื่อมต่อ WebSocket อยู่ จะนำพิกัดใหม่ไปขยับรถบนแผนที่ (Tracking) ได้ทันทีแบบไม่มีดีเลย์
5.  หากคนขับแจ้งปัญหา (Report Issue) ข้อมูลจะวิ่งเข้า API และถูก Broadcast ผ่าน WebSocket เพื่อให้ Web แสดง**กระดิ่งแจ้งเตือน** ทันที
