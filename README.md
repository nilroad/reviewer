# CSV to Voucher CSV Converters

This collection of Python scripts converts CSV files to voucher format for the club microservice. It supports both category-based and brand-based voucher generation by fetching product barcodes from the MySQL database.

## Features

- **Database Integration**: Connects to MySQL database using environment variables
- **Category-based Conversion**: Uses recursive CTE queries to find all child categories (similar to `GetCategoryChildren`)
- **Brand-based Conversion**: Finds products by brand ID and generates dynamic voucher titles
- **Product Barcode Extraction**: Fetches all product barcodes from category/brand hierarchies
- **CSV Format Transformation**: Converts from input formats to voucher format with proper field mapping
- **Dynamic Content**: Brand vouchers include brand names in titles dynamically
- **Error Handling**: Comprehensive logging and error handling for database operations

## Requirements

- Python 3.7+
- MySQL database access
- Required Python packages:
  - `mysql-connector-python==8.2.0` - MySQL database connectivity
  - `python-dotenv==1.0.0` - Environment variable loading from .env files
  
Install with: `pip install -r requirements.txt`

## Environment Variables

The script automatically loads database configuration from a `.env` file using python-dotenv (same variables as the Go application):

```bash
DATABASE_MYSQL_HOST=localhost
DATABASE_MYSQL_PORT=3306
DATABASE_MYSQL_USER=your_username
DATABASE_MYSQL_PASSWORD=your_password
DATABASE_MYSQL_NAME=club
```

The script will look for these variables in the following order:
1. `.env` file in the current directory (recommended)
2. System environment variables
3. Default values (where applicable)

## Setup

### Virtual Environment (Recommended)

For systems with externally managed Python environments:

```bash
# Create virtual environment
python3 -m venv venv

# Activate virtual environment
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Set environment variables (update with your credentials)
source config.template
```

### Environment Variables

The script automatically loads environment variables from a `.env` file. Set up your database credentials:

```bash
# Copy and edit the .env example file
cp env.example .env
# Edit .env with your actual database credentials

# Alternative: Export directly to shell
export DATABASE_MYSQL_HOST=localhost
export DATABASE_MYSQL_PORT=3306
export DATABASE_MYSQL_USER=your_username
export DATABASE_MYSQL_PASSWORD=your_password
export DATABASE_MYSQL_NAME=club
```

**Example .env file:**
```env
DATABASE_MYSQL_HOST=localhost
DATABASE_MYSQL_PORT=3306
DATABASE_MYSQL_USER=club_user
DATABASE_MYSQL_PASSWORD=your_secure_password
DATABASE_MYSQL_NAME=club
TZ=Asia/Tehran
```

## Usage

### Category-based Conversion

```bash
# Method 1: Using the provided shell script (recommended)
./run_conversion.sh "Category(2).csv" "converted_vouchers.csv"

# Method 2: Manual setup with virtual environment
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python convert_category_csv.py "Category(2).csv" "converted_vouchers.csv"
deactivate

# Method 3: Direct execution (if dependencies are already installed)
python3 convert_category_csv.py "Category(2).csv" "converted_vouchers.csv"
```

### Brand-based Conversion

```bash
# Method 1: Using the provided shell script (recommended)
./run_brand_conversion.sh "brand-sample.csv" "converted_brands.csv"

# Method 2: Manual setup with virtual environment
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python convert_brand_csv.py "brand-sample.csv" "converted_brands.csv"
deactivate

# Method 3: Direct execution (if dependencies are already installed)
python3 convert_brand_csv.py "brand-sample.csv" "converted_brands.csv"
```

## Input Format

### Category-based CSV Format

The category conversion script expects CSV files with the following columns:
- `cellphones`: Phone numbers (used for targeting)
- `category`: Category name to search in database
- `discount_percent`: Discount percentage
- `maximum_expiration_day_number`: Maximum days until expiration
- `max_discount_amount`: Maximum discount amount
- `voucher_title`: Title for the voucher

### Brand-based CSV Format

The brand conversion script expects CSV files with the following columns:
- `cellphones`: Phone numbers (used for targeting)
- `brand`: Brand name to search in database
- `discount_percent`: Discount percentage
- `maximum_expiration_day_number`: Maximum days until expiration
- `max_discount_amount`: Maximum discount amount

## Output Format

The script generates CSV files with voucher format (matching Go struct exactly) including:
- `voucher_title`: Original voucher title
- `voucher_subtitle`: Default Persian subtitle
- `voucher_description`: Default Persian description
- `minimum_accepted_invoice`: Minimum invoice amount (default: 10000)
- `discount_percent`: From input
- `expiration_date`: Default expiration date (11/06/2025)
- `maximum_expiration_day_number`: From input
- `maximum_allowed_count`: Default count (1)
- `is_unmergeable`: Default (0)
- `voucher_category_id`: Default category ID (75)
- `max_discount_amount`: From input (commas removed)
- `barcodes`: Product barcodes joined with `|` separator
- `cellphones`: Phone numbers from input (for targeting)
- `voucher_rule_type_number`: Default rule type (1)

## How It Works

### Category-based Conversion

1. **Category Lookup**: For each category name in the input, the script:
   - Searches for exact title match in the `categories` table
   - Uses recursive CTE query to find all child categories
   - Collects all category IDs (parent + descendants)

2. **Product Extraction**: For collected category IDs:
   - Queries all active products from those categories
   - Extracts product barcodes and titles
   - Formats barcodes with `|` separator as required

3. **Data Transformation**: 
   - Maps input fields to output format
   - Adds default values for required fields
   - Handles proper encoding for Persian text

### Brand-based Conversion

1. **Brand Lookup**: For each brand name in the input, the script:
   - Searches for exact title match in the `brands` table
   - Gets the brand ID for product queries

2. **Product Extraction**: For the brand ID:
   - Queries all active products with that brand_id
   - Extracts product barcodes and titles
   - Formats barcodes with `|` separator as required

3. **Dynamic Title Generation**: 
   - Creates voucher title including the brand name
   - Format: `تخفیف 10 درصدی ویژه شما روی برند "نام برند"`
   - Maps input fields to output format with brand-specific defaults

## Database Schema Requirements

The scripts expect the following tables to exist:
- `categories`: With columns `id`, `title`, `parent_id`, `deleted_at`
- `brands`: With columns `id`, `title`, `created_at`, `updated_at`
- `products`: With columns `id`, `barcode`, `title`, `category_id`, `brand_id`, `status`, `deleted_at`

## Error Handling

- **Missing Categories/Brands**: Logs warnings and skips rows with unfound categories or brands
- **No Products**: Skips categories/brands with no associated products
- **Database Errors**: Provides detailed error messages and graceful failure
- **Empty Rows**: Automatically skips empty or invalid CSV rows
- **Invalid Data**: Handles malformed CSV data and missing required fields

## Logging

The script provides comprehensive logging including:
- Database connection status
- Category lookup results
- Product count per category
- Processing statistics
- Error details

## Examples

### Category-based Conversion

Input CSV:
```csv
cellphones,category,discount_percent,maximum_expiration_day_number,max_discount_amount,voucher_title
989000009685,بستنی,10,7,"1,500,000",تخفیف ویژه شما برای بستنی
```

Output CSV:
```csv
voucher_title,voucher_subtitle,voucher_description,minimum_accepted_invoice,discount_percent,expiration_date,maximum_expiration_day_number,maximum_allowed_count,is_unmergeable,voucher_category_id,max_discount_amount,barcodes,cellphones
تخفیف ویژه شما برای بستنی,تخفیف 10 درصدی تا سقف 150 هزار تومان,این کوپن به صورت اختصاصی برای شما ساخته شده است,10000,10,06/11/2025,7,1,0,75,1500000,6260111310106|6260111310107,989000009685
```

### Brand-based Conversion

Input CSV:
```csv
cellphones,brand,discount_percent,maximum_expiration_day_number,max_discount_amount
989001215201,زر,10,7,"1,500,000"
```

Output CSV:
```csv
voucher_title,voucher_subtitle,voucher_description,minimum_accepted_invoice,discount_percent,expiration_date,maximum_expiration_day_number,maximum_allowed_count,is_unmergeable,voucher_category_id,max_discount_amount,barcodes,cellphones
تخفیف 10 درصدی ویژه شما روی برند "زر",تخفیف 10 درصدی تا سقف 150 هزار تومان,این کوپن به صورت اختصاصی برای شما ساخته شده است,10000,10,2025/11/06,7,1,0,75,1500000,6260111310108|6260111310109,989001215201
```

## Troubleshooting

1. **Database Connection Issues**: 
   - Verify environment variables are set correctly
   - Check database server accessibility
   - Ensure user has proper permissions

2. **Category/Brand Not Found**:
   - Verify category/brand names match exactly (case-sensitive)
   - Check if categories/brands exist in database
   - Look for extra spaces or special characters

3. **No Products Found**:
   - Verify products exist for the category/brand
   - Check product status is 'ACTIVE'
   - Ensure products have valid barcodes
   - For categories: Check if child categories have products
   - For brands: Verify brand_id is correctly set in products table

4. **Encoding Issues**:
   - Ensure input CSV is UTF-8 encoded
   - Check Persian text rendering
