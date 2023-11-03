package main

import (
	"bufio"
	"fmt"
	"github.com/esirangelomub/go-chat-application/configs"
	dbutils "github.com/esirangelomub/go-chat-application/database"
	"github.com/esirangelomub/go-chat-application/internal/entity"
	"os"
)

func main() {
	// load configs
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// database connection
	db, err := dbutils.InitializeDB(config)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{}, &entity.Chatroom{}, &entity.Message{})

	// Confirm action
	fmt.Println("This action will modify the database. Are you sure you want to continue? (yes/no)")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := scanner.Text()
	if response != "yes" {
		fmt.Println("Operation canceled.")
		return
	}

	// Execute SQL script
	fmt.Println("Adjusting database...")
	err = db.Exec("DO $$\n    BEGIN\n        IF EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_messages_chatroom_user') THEN\n            ALTER TABLE public.messages DROP CONSTRAINT fk_messages_chatroom_user;\n        END IF;\n        IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'chatroom_users') THEN\n            DROP TABLE public.chatroom_users CASCADE;\n        END IF;\n        IF EXISTS (SELECT FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'messages' AND column_name = 'chatroom_user_id') THEN\n            ALTER TABLE public.messages DROP COLUMN chatroom_user_id;\n        END IF;\n        IF EXISTS (SELECT FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'messages' AND column_name = 'timestamp') THEN\n            ALTER TABLE public.messages DROP COLUMN timestamp;\n        END IF;\n    END $$;").Error
	if err != nil {
		fmt.Println("Error adjusting database.")
		panic(err)
	}
	fmt.Println("Database adjusted successfully.")

	fmt.Println("Truncating database...")
	err = db.Exec("TRUNCATE public.users, public.chatrooms, public.messages RESTART IDENTITY CASCADE;").Error
	if err != nil {
		fmt.Println("Error truncating database.")
		panic(err)
	}
	fmt.Println("Database truncated successfully.")

	fmt.Println("Create users and chatroom's...")
	err = db.Exec("INSERT INTO public.users (id, name, email, password, type)\nVALUES\n    ('1bc8c7b5-4047-49e8-9ac6-f36a49510f2d', 'User1', 'user1@example.com', 'password1', 'USER'),\n    ('8b95244a-94e0-4097-adca-5d7b21ceb5d0', 'User2', 'user2@example.com', 'password2', 'USER'),\n    ('9ac7b8d0-5e2a-4f7b-8f6a-1e4f2c9a1b2e', 'User3', 'user3@example.com', 'password3', 'USER');\nINSERT INTO public.chatrooms (id, name, description, created_at)\nVALUES\n    ('bee6104c-bc5c-4837-87b5-9ce56c601ff0', 'Chatroom1', 'Description1', NOW()),\n    ('e4eaa663-bd7e-4a6a-97f7-cf8e6e7f8eb5', 'Chatroom2', 'Description2', NOW());").Error
	if err != nil {
		fmt.Println("Error creating users and chatroom's.")
		panic(err)
	}
	fmt.Println("Users and chatroom's created successfully.")
	fmt.Println("Create messages...")
	err = db.Exec("DO $$\n    DECLARE\n        i INT;\n        user_id UUID;\n        content_text TEXT;\n        created_at TIMESTAMP;\n    BEGIN\n        FOR i IN 1..200 LOOP\n                IF i % 3 = 0 THEN\n                    user_id := '1bc8c7b5-4047-49e8-9ac6-f36a49510f2d';\n                    content_text := 'Message from User1';\n                ELSIF i % 3 = 1 THEN\n                    user_id := '8b95244a-94e0-4097-adca-5d7b21ceb5d0';\n                    content_text := 'Message from User2';\n                ELSE\n                    user_id := '9ac7b8d0-5e2a-4f7b-8f6a-1e4f2c9a1b2e';\n                    content_text := 'Message from User3';\n                END IF;\n                IF i <= 50 THEN\n                    created_at := NOW();\n                ELSIF i <= 100 THEN\n                    created_at := NOW() - INTERVAL '1 day';\n                ELSIF i <= 150 THEN\n                    created_at := NOW() - INTERVAL '2 days';\n                ELSE\n                    created_at := NOW() - INTERVAL '1 month' - INTERVAL '1 day';\n                END IF;\n                IF i <= 100 THEN\n                    INSERT INTO public.messages (id, content, created_at, chatroom_id, user_id)\n                    VALUES (uuid_generate_v4(), content_text, created_at, 'bee6104c-bc5c-4837-87b5-9ce56c601ff0', user_id);\n                ELSE\n                    INSERT INTO public.messages (id, content, created_at, chatroom_id, user_id)\n                    VALUES (uuid_generate_v4(), content_text, created_at, 'e4eaa663-bd7e-4a6a-97f7-cf8e6e7f8eb5', user_id);\n                END IF;\n            END LOOP;\n    END $$;").Error
	if err != nil {
		fmt.Println("Error creating messages.")
		panic(err)
	}
	fmt.Println("Messages created successfully.")

	fmt.Println("Database setup completed successfully.")
}
