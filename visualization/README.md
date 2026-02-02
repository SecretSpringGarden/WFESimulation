# Workforce AI Transition Simulator - Visualization Tools

This directory contains Python scripts for visualizing simulation results from the Workforce AI Transition Simulator.

## Prerequisites

Install the required Python packages:

```bash
pip install -r requirements.txt
```

Or install individually:

```bash
pip install matplotlib pandas seaborn numpy plotly
```

## Scripts

### 1. plot_simulation.py

Creates static plots for single simulation results.

**Usage:**
```bash
# Basic usage
python plot_simulation.py simulation_report.json

# Specify output directory
python plot_simulation.py simulation_report.json -o plots/

# Generate only the dashboard
python plot_simulation.py simulation_report.json --dashboard-only

# Works with CSV files too
python plot_simulation.py simulation_report.csv
```

**Generated Plots:**
- `workforce_composition.png` - Human vs AI worker counts over time
- `revenue_productivity.png` - Revenue output and productivity trends
- `cost_analysis.png` - Cost utilization and budget analysis
- `equilibrium_analysis.png` - Final state analysis and key metrics
- `simulation_dashboard.png` - Comprehensive overview dashboard

### 2. plot_sensitivity.py

Creates visualizations for sensitivity analysis results.

**Usage:**
```bash
# Basic usage
python plot_sensitivity.py sensitivity_report.json

# Specify output directory
python plot_sensitivity.py sensitivity_report.json -o sensitivity_plots/

# Generate only the dashboard
python plot_sensitivity.py sensitivity_report.json --dashboard-only

# Works with detailed CSV files
python plot_sensitivity.py sensitivity_detailed.csv
```

**Generated Plots:**
- `parameter_rankings.png` - Parameter impact rankings
- `parameter_variations.png` - How metrics vary with parameter changes
- `sensitivity_heatmap.png` - Correlation heatmap (CSV data only)
- `outcome_distributions.png` - Distribution of key outcomes
- `sensitivity_dashboard.png` - Comprehensive sensitivity analysis overview

### 3. interactive_dashboard.py

Creates interactive web-based dashboards using Plotly.

**Usage:**
```bash
# Create simulation dashboard
python interactive_dashboard.py simulation_report.json

# Create sensitivity analysis dashboard
python interactive_dashboard.py sensitivity_report.json --sensitivity

# Specify output file and open in browser
python interactive_dashboard.py simulation_report.json -o dashboard.html --open
```

**Features:**
- Interactive plots with zoom, pan, and hover details
- Real-time data exploration
- Professional web-based interface
- Exportable as standalone HTML files

## Example Workflows

### Visualizing a Single Simulation

```bash
# Run simulation
./simulator -config examples/small_team_natural_attrition.json

# Create static plots
python visualization/plot_simulation.py simulation_report_*.json -o plots/

# Create interactive dashboard
python visualization/interactive_dashboard.py simulation_report_*.json --open
```

### Analyzing Sensitivity Results

```bash
# Run sensitivity analysis
./simulator -sensitivity -config examples/medium_team_fast_learning.yaml

# Create sensitivity plots
python visualization/plot_sensitivity.py sensitivity_report_*.json -o sensitivity_plots/

# Create interactive sensitivity dashboard
python visualization/interactive_dashboard.py sensitivity_report_*.json --sensitivity --open
```

### Batch Processing Multiple Results

```bash
# Process all simulation results in current directory
for file in simulation_report_*.json; do
    python visualization/plot_simulation.py "$file" -o "plots/$(basename "$file" .json)/"
done

# Process all sensitivity results
for file in sensitivity_report_*.json; do
    python visualization/plot_sensitivity.py "$file" -o "sensitivity_plots/$(basename "$file" .json)/"
done
```

## Output Formats

### Static Plots
- **Format**: PNG (300 DPI)
- **Size**: Optimized for presentations and reports
- **Style**: Professional seaborn styling with clear legends and labels

### Interactive Dashboards
- **Format**: HTML with embedded JavaScript
- **Features**: Zoom, pan, hover tooltips, legend toggling
- **Compatibility**: Works in any modern web browser
- **Sharing**: Self-contained files that can be shared easily

## Customization

### Modifying Plot Styles

Edit the style settings at the top of each script:

```python
# Change color palette
sns.set_palette("viridis")  # or "husl", "Set1", etc.

# Change plot style
plt.style.use('seaborn-v0_8')  # or 'ggplot', 'classic', etc.
```

### Adding New Visualizations

To add new plot types:

1. Create a new function following the existing pattern:
```python
def plot_new_analysis(df, output_dir):
    """Create new analysis plot."""
    fig, ax = plt.subplots(figsize=(12, 8))
    # Your plotting code here
    plt.savefig(output_dir / 'new_analysis.png', dpi=300, bbox_inches='tight')
    plt.close()
```

2. Call it from the main function:
```python
print("Creating new analysis...")
plot_new_analysis(df, output_dir)
```

### Handling Different Data Formats

The scripts automatically detect JSON vs CSV formats and adapt accordingly. For custom data formats:

1. Modify the `load_simulation_data()` function
2. Update column name mappings as needed
3. Add error handling for missing columns

## Troubleshooting

### Common Issues

1. **Missing Dependencies**
   ```bash
   pip install --upgrade matplotlib pandas seaborn numpy plotly
   ```

2. **Memory Issues with Large Datasets**
   - Use `--dashboard-only` flag for large simulations
   - Consider sampling data for very long time series

3. **Display Issues on Headless Systems**
   ```bash
   export MPLBACKEND=Agg  # Use non-interactive backend
   ```

4. **Font Rendering Issues**
   - Install system fonts: `sudo apt-get install fonts-dejavu-core`
   - Clear matplotlib cache: `rm -rf ~/.cache/matplotlib`

### Performance Tips

- Use `--dashboard-only` for quick overviews
- Process CSV files instead of JSON for large datasets
- Use batch processing scripts for multiple files
- Consider downsampling very long time series

## Integration with Other Tools

### Jupyter Notebooks

```python
import sys
sys.path.append('visualization/')
from plot_simulation import load_simulation_data, plot_workforce_composition

data, df = load_simulation_data('simulation_report.json')
plot_workforce_composition(df, Path('.'))
```

### Automated Reporting

Create automated reports by combining with shell scripts:

```bash
#!/bin/bash
# run_analysis.sh

# Run simulation
./simulator -config "$1"

# Get latest results
REPORT=$(ls -t simulation_report_*.json | head -1)

# Generate visualizations
python visualization/plot_simulation.py "$REPORT" -o "results/$(date +%Y%m%d)/"
python visualization/interactive_dashboard.py "$REPORT" -o "results/dashboard_$(date +%Y%m%d).html"

echo "Analysis complete. Results in results/ directory."
```

## Contributing

To contribute new visualization features:

1. Follow the existing code structure and naming conventions
2. Add appropriate error handling and user feedback
3. Update this README with new features
4. Test with various data formats and edge cases
5. Ensure plots are accessible (good contrast, clear labels)