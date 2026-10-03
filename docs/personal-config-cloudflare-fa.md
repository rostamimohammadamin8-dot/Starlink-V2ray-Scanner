<div dir="rtl">

# ساخت کانفیگ شخصی رایگان روی Cloudflare (بدون VPS و بدون کارت بانکی)

این راهنما با استفاده از پروژهٔ متن‌باز [BPB Panel](https://github.com/bia-pain-bache/BPB-Worker-Panel) یک لینک اشتراک VLESS / Trojan / Warp شخصی روی حساب رایگان Cloudflare خودتان می‌سازد. فقط یک ایمیل لازم است.

## محدودیت‌های واقعی (قبل از شروع بخوانید)

- **سقف روزانه:** هر Worker حدود ۱۰۰ هزار درخواست در روز برای VLESS/Trojan دارد؛ یعنی مناسب ۲ تا ۳ کاربر. کانفیگ‌های Warp این سقف را ندارند.
- **UDP:** روی Workers، VLESS و Trojan از UDP درست پشتیبانی نمی‌کنند؛ برای بازی‌هایی که UDP لازم دارند از کانفیگ‌های Warp پنل یا فهرست `subs/gaming.txt` همین پروژه استفاده کنید.
- **سرعت و پینگ:** به IP تمیز Cloudflare که از اینترنت خودتان پیدا می‌کنید بستگی دارد، نه به خود پنل.
- **قوانین Cloudflare:** استفاده از Workers به‌عنوان پروکسی با قوانین Cloudflare سازگار نیست و ممکن است حساب بسته شود. هیچ روشی این را تضمین نمی‌کند.

## مراحل نصب

۱. یک حساب جدید و **مخصوص همین کار** در Cloudflare بسازید: <https://dash.cloudflare.com/sign-up/> و ایمیل را تأیید کنید.

۲. پنل را با BPB Wizard نصب کنید (یکی از دو روش):

- نسخهٔ وب (ساده‌ترین): <https://wizard.bpb-panel.workers.dev> را باز کنید، طبق مراحل صفحه یک Token بسازید و نصب را بزنید.
- نسخهٔ ترمینال:

```bash
# Android (Termux) / Linux / macOS
bash <(curl -fsSL https://raw.githubusercontent.com/bia-pain-bache/BPB-Wizard/main/install.sh)
```

```powershell
# Windows PowerShell
irm https://raw.githubusercontent.com/bia-pain-bache/BPB-Wizard/main/install.ps1 | iex
```

۳. بعد از نصب، آدرس پنل و رمز آن را جایی امن نگه دارید و وارد پنل شوید.

۴. از بخش Subscription لینک مخصوص کلاینت خود (v2rayNG، v2rayN، Streisand، Sing-box، Clash) را کپی کنید و در برنامه Import کنید.

راهنمای کامل تنظیمات و استفاده (فارسی): <https://bia-pain-bache.github.io/BPB-Worker-Panel/fa/>

## کم کردن ریسک مسدود شدن و برگشت سریع

- لینک اشتراک را **عمومی پخش نکنید**؛ فقط برای خودتان و چند نفر نزدیک.
- برای این کار از حساب Cloudflare جداگانه استفاده کنید تا اگر بسته شد به چیز دیگری آسیب نخورد.
- از تنظیمات پنل نسخهٔ پشتیبان بگیرید؛ با BPB Wizard نصب مجدد روی یک حساب تازه چند دقیقه طول می‌کشد.
- همیشه یک جایگزین داشته باشید: فهرست‌های تست‌شدهٔ همین پروژه (`cleaned_configs.txt` و پوشهٔ `subs/`) هر ۶ ساعت به‌روز می‌شوند.

</div>
