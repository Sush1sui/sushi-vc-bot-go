## sushi-vc-bot-go

Discord voice-channel bot with configurable role and channel exceptions.

## Environment Variables

Required:

- `BOT_TOKEN`
- `APP_ID`
- `MONGODB_URI`
- `MONGODB_NAME`
- `CATEGORY_JTC_COLLECTION_NAME`
- `CUSTOM_VC_COLLECTION_NAME`

Optional:

- `PORT` (default: `7694`)
- `SERVER_URL`
- `GUILD_SETTINGS_COLLECTION_NAME` (default: `guild_settings`)

Notes:

- Role/channel exceptions are loaded from MongoDB per guild.
- If no document exists for a guild, all lists are treated as empty until you configure them via `/guild-settings`.

### Guild Settings Document

Collection: `guild_settings` (or value from `GUILD_SETTINGS_COLLECTION_NAME`)

```json
{
  "guild_id": "123456789012345678",
  "voice_exception_role_ids": ["1292473360114122784", "1310186525606154340"],
  "protected_role_ids": ["1292473360114122784", "1310186525606154340"],
  "ignored_channel_ids": ["140000000000000000"]
}
```

## Slash Commands For Settings

Use these admin commands in Discord:

- `/guild-settings create target:<voice_exception_roles|protected_roles|ignored_channels> ids:<comma-separated-ids>`
- `/guild-settings update target:<voice_exception_roles|protected_roles|ignored_channels> ids:<comma-separated-ids>`
- `/guild-settings delete target:<voice_exception_roles|protected_roles|ignored_channels> ids:<comma-separated-ids>`
- `/guild-settings delete target:<voice_exception_roles|protected_roles|ignored_channels>` (clears that entire list)
- `/guild-settings view`
- `/guild-settings reset` (deletes guild settings document)

Examples:

- `/guild-settings create target:voice_exception_roles ids:1292473360114122784,1310186525606154340`
- `/guild-settings update target:protected_roles ids:1292473360114122784,1310186525606154340`
- `/guild-settings delete target:protected_roles ids:1310186525606154340`
