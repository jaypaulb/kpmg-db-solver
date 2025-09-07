# Database Schema Extraction

This directory contains scripts to extract the PostgreSQL database schema from a Canvus server installation.

## Files

- `extract-schema.ps1` - PowerShell script for schema extraction
- `extract-schema.bat` - Windows batch file wrapper
- `README-schema.md` - This documentation

## Prerequisites

1. **PostgreSQL Client Tools** - The `psql` command must be available in PATH
   - Download from: https://www.postgresql.org/download/
   - Or install via package manager: `choco install postgresql` (Chocolatey)

2. **Canvus Server Installation** - The `mt-canvus-server.ini` file must exist at:
   ```
   C:\ProgramData\MultiTaction\canvus\mt-canvus-server.ini
   ```

3. **Database Access** - The script will use the credentials from the INI file to connect to the database

## Usage

### Option 1: PowerShell (Recommended)
```powershell
.\scripts\extract-schema.ps1
```

### Option 2: Batch File (Windows)
```cmd
.\scripts\extract-schema.bat
```

### Option 3: Custom Parameters
```powershell
.\scripts\extract-schema.ps1 -IniPath "C:\Custom\Path\mt-canvus-server.ini" -OutputPath "custom-schema.sql"
```

## What It Does

The script will:

1. **Read Configuration** - Parse the `mt-canvus-server.ini` file
2. **Extract Database Details** - Get database name, username, password, and port
3. **Connect to Database** - Use `psql` to connect to the PostgreSQL database
4. **Export Schema** - Generate a comprehensive schema file including:
   - Database schemas
   - Table structures
   - Indexes
   - Sequences
   - Functions
   - Views
   - Foreign key relationships

## Output

The script generates a file called `canvus-schema.sql` containing:
- Complete database structure
- Table definitions with column types and constraints
- Index information
- Foreign key relationships
- Function and view definitions

## Example Output

```
🔍 Extracting Canvus Database Schema
=====================================
📄 Reading configuration from: C:\ProgramData\MultiTaction\canvus\mt-canvus-server.ini
📊 Database Configuration:
   Database: mt_canvus
   Username: mt_canvus
   Port: 5432
   Password: [HIDDEN]
✅ Found psql at: C:\Program Files\PostgreSQL\14\bin\psql.exe
🔗 Connecting to database...
✅ Successfully connected to database
📝 Extracting schema to: canvus-schema.sql
✅ Schema successfully extracted to: canvus-schema.sql
📊 Export Statistics:
   File Size: 45.67 KB
   Lines: 1,234

🎉 Schema extraction completed successfully!
You can now share the schema file: canvus-schema.sql
```

## Troubleshooting

### Common Issues

1. **"psql command not found"**
   - Install PostgreSQL client tools
   - Ensure `psql` is in your system PATH

2. **"INI file not found"**
   - Verify Canvus server is installed
   - Check the path: `C:\ProgramData\MultiTaction\canvus\mt-canvus-server.ini`

3. **"Failed to connect to database"**
   - Verify database credentials in the INI file
   - Ensure PostgreSQL service is running
   - Check firewall settings

4. **PowerShell Execution Policy**
   - Run: `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser`
   - Or use the batch file wrapper

## Security Notes

- The script uses the `PGPASSWORD` environment variable for password authentication
- The environment variable is cleaned up after execution
- Passwords are not displayed in the output (shown as [HIDDEN])

## Next Steps

After extracting the schema:

1. **Review the schema file** - Understand the database structure
2. **Share for analysis** - Provide the schema file for database analysis
3. **Use in development** - Reference the schema for application development
4. **Document relationships** - Map out table relationships and dependencies
