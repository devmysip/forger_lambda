package utilis

import "forger/gita/models"

func GetNotificationTemplates() []models.NotificationModel {
	return []models.NotificationModel{
		{
			Title: "Shuruat ka Pehla Kadam!",
			Body:  "Shuruat karo! Aaj ek adhyay padhkar apne din ko pavitra banao. 🕉️",
		},
		{
			Title: "Aaj Ka Adhyay: Shanti aur Gyaan Ki Ore",
			Body:  "Kya aapne aaj ka adhyay padha? Thoda samay nikaalo aur Gita se shanti pao. 🙏",
		},
		{
			Title: "Jeevan Ka Satya: Adhyay Ka Sandesh",
			Body:  "Bhagavad Gita padhna shuru karo aur jeevan ke sach ko samjho. Aaj ka adhyay zarur dekhein! 📖",
		},
		{
			Title: "5 Minute Mein Gyaan Ka Sagar!",
			Body:  "Sirf 5 minute ka samay dekar apne vicharon ko saf karo. Aaj ka adhyay padhne ka samay hai! ⏳",
		},
		{
			Title: "Krishna Ka Sandesh Tumhare Liye",
			Body:  "Krishna ke gyaan se apne din ko roshan karo. Bhagavad Gita ka adhyay tumhara intezaar kar raha hai! 🌟",
		},
		{
			Title: "Gyaan Ka Naya Dwaar Kholiye",
			Body:  "Har adhyay ek naye gyaan ka dwar kholta hai. Aaj ka adhyay zarur padho. 🚪📘",
		},
		{
			Title: "Safar Ka Saptah: Antim Kadam, Naya Gyaan",
			Body:  "Ab tak ke safar ka anand lo! Aaj ka adhyay padhkar apne Gyaan ka vruddhi karo. 🕉️✨",
		},
	}
}
