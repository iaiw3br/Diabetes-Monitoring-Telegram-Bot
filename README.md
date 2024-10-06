This project is designed to help monitor and manage critical diabetes-related data for a friend’s daughter who has diabetes. The application integrates with the Nightscout API to fetch real-time glucose data and sends notifications to Telegram in case of critical glucose levels. Additionally, it allows temporary message pauses when insulin injections are administered, ensuring only relevant and timely notifications are sent.

Features

	•	Critical Value Alerts: Automatically sends notifications to a Telegram chat when glucose levels fall outside safe ranges, ensuring immediate action when needed.
	•	Pause Notifications: Pauses alerts when insulin injections are administered to avoid unnecessary alerts during treatment periods.

Technologies Used

	•	Golang: Core programming language used to build the application.
	•	Telegram API: Used to send real-time alerts and updates directly to a designated Telegram chat.
	•	Nightscout API (nightscout-jino.ru/api): Integrated to fetch real-time glucose data for monitoring.

How It Works

	1.	Fetch Data: The application periodically checks glucose data from the Nightscout API.
	2.	Send Alerts: If glucose levels reach critical values, a message is sent to a Telegram chat.
	3.	Pause Feature: Allows pausing of notifications when insulin injections are being administered to avoid unnecessary interruptions.

How to Use

	1.	Set up your Nightscout instance or use an existing API endpoint.
	2.	Connect the bot to your Telegram account using the Telegram API.
	3.	Customize the critical value thresholds and pause settings to match the patient’s treatment plan.

This project provides a simple, automated solution for monitoring glucose levels and ensuring timely action through real-time notifications in Telegram.
