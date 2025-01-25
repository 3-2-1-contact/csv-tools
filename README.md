## 📖 Usage

csv-tools provides two main commands: `merge` and `split` for manipulating CSV files.

### Merge Command

The merge command combines multiple CSV files into a single file, with the option to add a new column containing the source filename.

```bash
csvtools merge [flags]

Flags:
  --in string    Input CSV files (comma-separated) [required]
  --out string   Output CSV file (optional, outputs to stdout if not specified)
  --col string   Name of new column to add with filename values
```

**Examples:**
```bash
# Merge county data files and add a 'county' column
csvtools merge --col county --in "comal.csv,hays.csv" --out all.csv

# Merge state data files and pipe to split command
csvtools merge --col state --in "texas.csv,florida.csv" | csvtools split --filter "category=Aircraft" --split state
```

### Split Command

The split command divides a CSV file into multiple files based on specified criteria and can filter rows based on column values.

```bash
csvtools split [flags]

Flags:
  --in string      Input CSV file (optional, reads from stdin if not provided)
  --filter string  Filter in format column=value1,value2,... [required]
  --split string   Column to split by [required]
  --format string  Column format specifications
```

**Format Specifications:**
- Available formats: string, int, decimal:<precision>
- Format string example: "Total Value:decimal:6,'Account Number':int"

**Examples:**
```bash
# Split a CSV file by state column, filtering for Aircraft category
csvtools split --in combined.csv --filter "category=Aircraft" --split state

# Split with column formatting
csvtools split --in financial.csv --filter "type=transaction" --split account --format "Total:decimal:2,ID:int"
```

### Pipeline Usage

Commands can be chained together using Unix-style pipes:

```bash
# Merge files, then split the result
csvtools merge --col region --in "north.csv,south.csv" | csvtools split --filter "type=sales" --split region

# Process and output to file
csvtools merge --in "*.csv" --col source | csvtools split --filter "status=active" --split department > output.csv
```

### Common Patterns

1. **Combining Regional Data:**
```bash
csvtools merge --col region --in "region*.csv" --out combined.csv
```

2. **Filtering and Splitting Large Datasets:**
```bash
csvtools split --in massive.csv --filter "year=2024" --split department
```

3. **Format Numbers in Financial Data:**
```bash
csvtools split --in transactions.csv --filter "type=expense" --split category --format "Amount:decimal:2"
```

## 💡 Tips

- Use `--out` flag to save to a file, or omit it to output to stdout
- When reading from stdin, omit the `--in` flag
- File patterns (e.g., `*.csv`) can be used with the `--in` flag
- Format specifications can handle spaces in column names using single quotes
- Multiple filter values can be comma-separated
---