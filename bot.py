import logging

import telebot

TOKEN_PATH = "./config/token.txt"
ADMIN_ID = 852084868
XOR_KEY = 19011996

WELCOME_TEXT = (
    "привет!\n"
    "это бот канала «[давай на ты](https://t.me/+eFiRVo8U-NhlNjEy)», "
    "и всё, что ты сюда напишешь — анонимно отправится Ксюше 📖\n\n"
    "она уже ждёт твои вопросы и истории!"
)


def read_token(path: str) -> str:
    with open(path, "r", encoding="utf-8") as f:
        return f.readline().strip()


def encrypt_id(chat_id: int) -> int:
    return chat_id ^ XOR_KEY


def decrypt_id(encrypted: int) -> int:
    return encrypted ^ XOR_KEY


logging.basicConfig(
    level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s"
)

bot = telebot.TeleBot(read_token(TOKEN_PATH))
bot.remove_webhook()
logging.info("Authorized on account %s", bot.get_me().username)


@bot.message_handler(commands=["start"])
def handle_start(message):
    bot.send_message(
        message.chat.id,
        WELCOME_TEXT,
        parse_mode="Markdown",
        disable_web_page_preview=True,
    )


@bot.message_handler(
    func=lambda m: m.chat.type == "private" and m.chat.id != ADMIN_ID,
    content_types=[
        "text",
        "photo",
        "audio",
        "document",
        "video",
        "voice",
        "video_note",
    ],
)
def forward_to_admin(message):
    chat_id = str(encrypt_id(message.chat.id))
    caption_prefix = f"{chat_id} :\n"
    caption = caption_prefix + (message.caption or "")

    if message.content_type == "text":
        if not message.text:
            return
        bot.send_message(ADMIN_ID, f"{chat_id} :\n{message.text}")
    elif message.content_type == "photo":
        bot.send_photo(ADMIN_ID, message.photo[-1].file_id, caption=caption)
    elif message.content_type == "audio":
        bot.send_audio(ADMIN_ID, message.audio.file_id, caption=caption)
    elif message.content_type == "document":
        bot.send_document(ADMIN_ID, message.document.file_id, caption=caption)
    elif message.content_type == "video":
        bot.send_video(ADMIN_ID, message.video.file_id, caption=caption)
    elif message.content_type == "voice":
        bot.send_message(ADMIN_ID, f"{chat_id} :")
        bot.send_voice(ADMIN_ID, message.voice.file_id)
    elif message.content_type == "video_note":
        bot.send_message(ADMIN_ID, f"{chat_id} :")
        bot.send_video_note(ADMIN_ID, message.video_note.file_id)


@bot.message_handler(
    func=lambda m: m.chat.id == ADMIN_ID and m.reply_to_message is not None,
    content_types=[
        "text",
        "photo",
        "audio",
        "document",
        "video",
        "voice",
        "video_note",
    ],
)
def admin_reply(message):
    original = message.reply_to_message
    text_source = original.caption if original.caption else (original.text or "")
    words = text_source.split()
    if not words:
        return
    try:
        encrypted = int(words[0])
    except ValueError:
        return
    target_chat = decrypt_id(encrypted)

    if message.content_type == "text":
        bot.send_message(target_chat, message.text)
    elif message.content_type == "photo":
        bot.send_photo(
            target_chat, message.photo[-1].file_id, caption=message.caption
        )
    elif message.content_type == "audio":
        bot.send_audio(target_chat, message.audio.file_id, caption=message.caption)
    elif message.content_type == "document":
        bot.send_document(
            target_chat, message.document.file_id, caption=message.caption
        )
    elif message.content_type == "video":
        bot.send_video(target_chat, message.video.file_id, caption=message.caption)
    elif message.content_type == "voice":
        bot.send_voice(target_chat, message.voice.file_id)
    elif message.content_type == "video_note":
        bot.send_video_note(target_chat, message.video_note.file_id)


if __name__ == "__main__":
    bot.infinity_polling(timeout=60, long_polling_timeout=60)
