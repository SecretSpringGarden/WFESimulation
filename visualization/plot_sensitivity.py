#!/usr/bin/env python3
"""
Workforce AI Transition Simulator - Sensitivity Analysis Visualization

This script generates visualizations from sensitivity analysis results including:
- Parameter impact rankings
- Time to equilibrium variations
- Workforce composition variations
- Heatmaps and correlation analysis

Usage:
    python plot_sensitivity.py sensitivity_report.json
    python plot_sensitivity.py sensitivity_detailed.csv
"""

import argparse
import json
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import numpy as np
from pathlib import Path
import sys

# Set style for better-looking plots
plt.style.use('seaborn-v0_8')
sns.set_palette("viridis")

def load_sensitivity_data(file_path):
    """Load sensitivity analysis data from JSON or CSV file."""
    file_path = Path(file_path)
    
    if file_path.suffix.lower() == '.json':
        with open(file_path, 'r') as f:
            data = json.load(f)
        return data, None
    
    elif file_path.suffix.lower() == '.csv':
        df = pd.read_csv(file_path)
        return None, df
    
    else:
        print(f"Error: Unsupported file format {file_path.suffix}")
        return None, None

def plot_parameter_rankings(data, output_dir):
    """Plot parameter impact rankings."""
    if not data or 'ParameterRankings' not in data:
        print("Warning: No parameter rankings found in data")
        return
    
    rankings = data['ParameterRankings']
    
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(16, 8))
    
    # Time to equilibrium impact
    time_impacts = []
    param_names = []
    for param in rankings:
        if 'TimeToEquilibriumImpact' in param:
            time_impacts.append(param['TimeToEquilibriumImpact'])
            param_names.append(param['ParameterName'])
    
    if time_impacts:
        y_pos = np.arange(len(param_names))
        bars1 = ax1.barh(y_pos, time_impacts, color='skyblue')
        ax1.set_yticks(y_pos)
        ax1.set_yticklabels(param_names)
        ax1.set_xlabel('Impact Score')
        ax1.set_title('Parameter Impact on Time to Equilibrium')
        ax1.grid(True, alpha=0.3)
        
        # Add value labels on bars
        for i, bar in enumerate(bars1):
            width = bar.get_width()
            ax1.text(width + 0.01, bar.get_y() + bar.get_height()/2, 
                    f'{width:.3f}', ha='left', va='center')
    
    # Workforce composition impact
    comp_impacts = []
    for param in rankings:
        if 'WorkforceCompositionImpact' in param:
            comp_impacts.append(param['WorkforceCompositionImpact'])
    
    if comp_impacts:
        bars2 = ax2.barh(y_pos, comp_impacts, color='lightcoral')
        ax2.set_yticks(y_pos)
        ax2.set_yticklabels(param_names)
        ax2.set_xlabel('Impact Score')
        ax2.set_title('Parameter Impact on Workforce Composition')
        ax2.grid(True, alpha=0.3)
        
        # Add value labels on bars
        for i, bar in enumerate(bars2):
            width = bar.get_width()
            ax2.text(width + 0.01, bar.get_y() + bar.get_height()/2, 
                    f'{width:.3f}', ha='left', va='center')
    
    plt.tight_layout()
    plt.savefig(output_dir / 'parameter_rankings.png', dpi=300, bbox_inches='tight')
    plt.close()

def plot_parameter_variations(data, output_dir):
    """Plot how key metrics vary with parameter changes."""
    if not data or 'SensitivityResults' not in data:
        print("Warning: No sensitivity results found in data")
        return
    
    results = data['SensitivityResults']
    
    # Create subplots for each parameter
    n_params = len(results)
    cols = min(3, n_params)
    rows = (n_params + cols - 1) // cols
    
    fig, axes = plt.subplots(rows, cols, figsize=(5*cols, 4*rows))
    if n_params == 1:
        axes = [axes]
    elif rows == 1:
        axes = [axes]
    else:
        axes = axes.flatten()
    
    for i, param_result in enumerate(results):
        if i >= len(axes):
            break
            
        ax = axes[i]
        param_name = param_result['ParameterName']
        
        if 'ParameterValues' in param_result and 'TimeToEquilibriumByValue' in param_result:
            param_values = param_result['ParameterValues']
            time_values = [param_result['TimeToEquilibriumByValue'].get(str(v), 0) 
                          for v in param_values]
            
            ax.plot(param_values, time_values, marker='o', linewidth=2, markersize=6)
            ax.set_xlabel(param_name)
            ax.set_ylabel('Time to Equilibrium')
            ax.set_title(f'Impact of {param_name}')
            ax.grid(True, alpha=0.3)
    
    # Hide unused subplots
    for i in range(n_params, len(axes)):
        axes[i].set_visible(False)
    
    plt.tight_layout()
    plt.savefig(output_dir / 'parameter_variations.png', dpi=300, bbox_inches='tight')
    plt.close()

def plot_sensitivity_heatmap(df, output_dir):
    """Create a heatmap of parameter sensitivity from CSV data."""
    if df is None:
        print("Warning: No CSV data available for heatmap")
        return
    
    # Identify parameter columns and outcome columns
    param_cols = []
    outcome_cols = []
    
    for col in df.columns:
        if col.startswith('Param_') or col in ['InitialHumans', 'FixedBudget', 'NaturalRate']:
            param_cols.append(col)
        elif col in ['TimeToEquilibrium', 'FinalHumans', 'FinalAIAgents', 'FinalRevenue']:
            outcome_cols.append(col)
    
    if not param_cols or not outcome_cols:
        print("Warning: Could not identify parameter and outcome columns for heatmap")
        return
    
    # Calculate correlations
    corr_data = []
    for param in param_cols:
        for outcome in outcome_cols:
            if param in df.columns and outcome in df.columns:
                corr = df[param].corr(df[outcome])
                corr_data.append({
                    'Parameter': param,
                    'Outcome': outcome,
                    'Correlation': corr
                })
    
    if not corr_data:
        return
    
    # Create correlation matrix
    corr_df = pd.DataFrame(corr_data)
    corr_matrix = corr_df.pivot(index='Parameter', columns='Outcome', values='Correlation')
    
    # Create heatmap
    plt.figure(figsize=(12, 8))
    sns.heatmap(corr_matrix, annot=True, cmap='RdBu_r', center=0, 
                square=True, fmt='.3f', cbar_kws={'label': 'Correlation'})
    plt.title('Parameter-Outcome Correlation Heatmap')
    plt.tight_layout()
    plt.savefig(output_dir / 'sensitivity_heatmap.png', dpi=300, bbox_inches='tight')
    plt.close()

def plot_parameter_distributions(df, output_dir):
    """Plot distributions of outcomes for different parameter values."""
    if df is None:
        return
    
    # Find key outcome columns
    outcome_cols = ['TimeToEquilibrium', 'FinalHumans', 'FinalAIAgents', 'FinalRevenue']
    available_outcomes = [col for col in outcome_cols if col in df.columns]
    
    if not available_outcomes:
        print("Warning: No outcome columns found for distribution plots")
        return
    
    n_outcomes = len(available_outcomes)
    fig, axes = plt.subplots(2, 2, figsize=(15, 12))
    axes = axes.flatten()
    
    for i, outcome in enumerate(available_outcomes[:4]):
        ax = axes[i]
        
        # Create histogram
        ax.hist(df[outcome], bins=20, alpha=0.7, color='skyblue', edgecolor='black')
        ax.set_xlabel(outcome)
        ax.set_ylabel('Frequency')
        ax.set_title(f'Distribution of {outcome}')
        ax.grid(True, alpha=0.3)
        
        # Add statistics
        mean_val = df[outcome].mean()
        std_val = df[outcome].std()
        ax.axvline(mean_val, color='red', linestyle='--', linewidth=2, label=f'Mean: {mean_val:.2f}')
        ax.axvline(mean_val + std_val, color='orange', linestyle=':', alpha=0.7, label=f'±1 Std: {std_val:.2f}')
        ax.axvline(mean_val - std_val, color='orange', linestyle=':', alpha=0.7)
        ax.legend()
    
    # Hide unused subplots
    for i in range(n_outcomes, 4):
        axes[i].set_visible(False)
    
    plt.tight_layout()
    plt.savefig(output_dir / 'outcome_distributions.png', dpi=300, bbox_inches='tight')
    plt.close()

def create_sensitivity_dashboard(data, df, output_dir):
    """Create a comprehensive sensitivity analysis dashboard."""
    fig = plt.figure(figsize=(20, 16))
    gs = fig.add_gridspec(4, 3, hspace=0.3, wspace=0.3)
    
    # Parameter rankings (top row)
    if data and 'ParameterRankings' in data:
        rankings = data['ParameterRankings']
        
        # Time to equilibrium impact
        ax1 = fig.add_subplot(gs[0, :2])
        time_impacts = []
        param_names = []
        for param in rankings:
            if 'TimeToEquilibriumImpact' in param:
                time_impacts.append(param['TimeToEquilibriumImpact'])
                param_names.append(param['ParameterName'])
        
        if time_impacts:
            y_pos = np.arange(len(param_names))
            ax1.barh(y_pos, time_impacts, color='skyblue')
            ax1.set_yticks(y_pos)
            ax1.set_yticklabels(param_names)
            ax1.set_title('Parameter Impact on Time to Equilibrium')
            ax1.grid(True, alpha=0.3)
        
        # Summary statistics
        ax2 = fig.add_subplot(gs[0, 2])
        if df is not None and 'TimeToEquilibrium' in df.columns:
            stats_text = f"""Sensitivity Analysis Summary
            
Total Runs: {len(df)}
Avg Time to Equilibrium: {df['TimeToEquilibrium'].mean():.1f}
Min Time: {df['TimeToEquilibrium'].min():.0f}
Max Time: {df['TimeToEquilibrium'].max():.0f}
Std Dev: {df['TimeToEquilibrium'].std():.1f}"""
            ax2.text(0.1, 0.5, stats_text, transform=ax2.transAxes, fontsize=10,
                    verticalalignment='center', fontfamily='monospace')
            ax2.set_title('Summary Statistics')
            ax2.axis('off')
    
    # Parameter variations (middle rows)
    if data and 'SensitivityResults' in data:
        results = data['SensitivityResults']
        
        for i, param_result in enumerate(results[:6]):  # Show up to 6 parameters
            row = 1 + i // 3
            col = i % 3
            if row >= 3:
                break
                
            ax = fig.add_subplot(gs[row, col])
            param_name = param_result['ParameterName']
            
            if 'ParameterValues' in param_result and 'TimeToEquilibriumByValue' in param_result:
                param_values = param_result['ParameterValues']
                time_values = [param_result['TimeToEquilibriumByValue'].get(str(v), 0) 
                              for v in param_values]
                
                ax.plot(param_values, time_values, marker='o', linewidth=2, markersize=4)
                ax.set_xlabel(param_name)
                ax.set_ylabel('Time to Equilibrium')
                ax.set_title(f'{param_name} Impact')
                ax.grid(True, alpha=0.3)
    
    # Outcome distributions (bottom row)
    if df is not None:
        outcome_cols = ['TimeToEquilibrium', 'FinalHumans', 'FinalAIAgents']
        available_outcomes = [col for col in outcome_cols if col in df.columns]
        
        for i, outcome in enumerate(available_outcomes[:3]):
            ax = fig.add_subplot(gs[3, i])
            ax.hist(df[outcome], bins=15, alpha=0.7, color='lightcoral', edgecolor='black')
            ax.set_xlabel(outcome)
            ax.set_ylabel('Frequency')
            ax.set_title(f'{outcome} Distribution')
            ax.grid(True, alpha=0.3)
    
    plt.suptitle('Sensitivity Analysis Dashboard', fontsize=20, weight='bold')
    plt.savefig(output_dir / 'sensitivity_dashboard.png', dpi=300, bbox_inches='tight')
    plt.close()

def main():
    parser = argparse.ArgumentParser(description='Visualize Workforce AI Transition Sensitivity Analysis')
    parser.add_argument('input_file', help='Path to sensitivity analysis results file (JSON or CSV)')
    parser.add_argument('-o', '--output', default='sensitivity_plots', help='Output directory for plots')
    parser.add_argument('--dashboard-only', action='store_true', help='Generate only the dashboard')
    
    args = parser.parse_args()
    
    # Create output directory
    output_dir = Path(args.output)
    output_dir.mkdir(exist_ok=True)
    
    # Load data
    print(f"Loading sensitivity analysis data from {args.input_file}...")
    data, df = load_sensitivity_data(args.input_file)
    
    if data is None and df is None:
        print("Error: Could not load sensitivity analysis data")
        sys.exit(1)
    
    if df is not None:
        print(f"Loaded {len(df)} sensitivity analysis runs")
    
    if args.dashboard_only:
        print("Creating sensitivity analysis dashboard...")
        create_sensitivity_dashboard(data, df, output_dir)
    else:
        # Generate all plots
        if data:
            print("Creating parameter ranking plots...")
            plot_parameter_rankings(data, output_dir)
            
            print("Creating parameter variation plots...")
            plot_parameter_variations(data, output_dir)
        
        if df is not None:
            print("Creating sensitivity heatmap...")
            plot_sensitivity_heatmap(df, output_dir)
            
            print("Creating outcome distribution plots...")
            plot_parameter_distributions(df, output_dir)
        
        print("Creating sensitivity analysis dashboard...")
        create_sensitivity_dashboard(data, df, output_dir)
    
    print(f"Plots saved to {output_dir}/")
    print("Generated files:")
    for plot_file in output_dir.glob('*.png'):
        print(f"  - {plot_file.name}")

if __name__ == '__main__':
    main()