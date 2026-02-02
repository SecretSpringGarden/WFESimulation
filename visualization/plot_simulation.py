#!/usr/bin/env python3
"""
Workforce AI Transition Simulator - Visualization Script

This script generates visualizations from simulation results including:
- Workforce composition over time
- Revenue output over time
- Cost and productivity trends
- Equilibrium analysis

Usage:
    python plot_simulation.py simulation_report.json
    python plot_simulation.py simulation_report.csv
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
sns.set_palette("husl")

def load_simulation_data(file_path):
    """Load simulation data from JSON or CSV file."""
    file_path = Path(file_path)
    
    if file_path.suffix.lower() == '.json':
        with open(file_path, 'r') as f:
            data = json.load(f)
        
        # Convert time series to DataFrame
        if 'TimeSeries' in data:
            df = pd.DataFrame(data['TimeSeries'])
            return data, df
        else:
            print("Error: No TimeSeries data found in JSON file")
            return None, None
    
    elif file_path.suffix.lower() == '.csv':
        df = pd.read_csv(file_path)
        return None, df
    
    else:
        print(f"Error: Unsupported file format {file_path.suffix}")
        return None, None

def plot_workforce_composition(df, output_dir):
    """Plot workforce composition over time."""
    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(12, 10))
    
    # Extract workforce data (assuming nested structure in JSON)
    if 'Workforce.Humans.Total' in df.columns:
        humans_col = 'Workforce.Humans.Total'
        ai_col = 'Workforce.AIAgents.Total'
    else:
        # Try alternative column names for CSV
        humans_col = 'TotalHumans'
        ai_col = 'TotalAIAgents'
        if humans_col not in df.columns:
            print("Warning: Could not find workforce composition columns")
            return
    
    # Plot absolute numbers
    ax1.plot(df['TimeStep'], df[humans_col], label='Human Workers', linewidth=2, marker='o', markersize=4)
    ax1.plot(df['TimeStep'], df[ai_col], label='AI Agents', linewidth=2, marker='s', markersize=4)
    ax1.set_xlabel('Time Step')
    ax1.set_ylabel('Number of Workers')
    ax1.set_title('Workforce Composition Over Time')
    ax1.legend()
    ax1.grid(True, alpha=0.3)
    
    # Plot as stacked area chart
    ax2.fill_between(df['TimeStep'], 0, df[humans_col], label='Human Workers', alpha=0.7)
    ax2.fill_between(df['TimeStep'], df[humans_col], df[humans_col] + df[ai_col], 
                     label='AI Agents', alpha=0.7)
    ax2.set_xlabel('Time Step')
    ax2.set_ylabel('Number of Workers')
    ax2.set_title('Workforce Composition (Stacked)')
    ax2.legend()
    ax2.grid(True, alpha=0.3)
    
    plt.tight_layout()
    plt.savefig(output_dir / 'workforce_composition.png', dpi=300, bbox_inches='tight')
    plt.close()

def plot_revenue_and_productivity(df, output_dir):
    """Plot revenue output and productivity over time."""
    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(12, 10))
    
    # Revenue plot
    if 'RevenueOutput' in df.columns:
        ax1.plot(df['TimeStep'], df['RevenueOutput'], label='Revenue Output', 
                linewidth=2, color='green', marker='o', markersize=4)
        ax1.set_xlabel('Time Step')
        ax1.set_ylabel('Revenue Output')
        ax1.set_title('Revenue Output Over Time')
        ax1.grid(True, alpha=0.3)
        ax1.ticklabel_format(style='plain', axis='y')
    
    # Productivity plot
    if 'TotalProductivity' in df.columns:
        ax2.plot(df['TimeStep'], df['TotalProductivity'], label='Total Productivity', 
                linewidth=2, color='blue', marker='s', markersize=4)
        ax2.set_xlabel('Time Step')
        ax2.set_ylabel('Total Productivity')
        ax2.set_title('Total Productivity Over Time')
        ax2.grid(True, alpha=0.3)
    
    plt.tight_layout()
    plt.savefig(output_dir / 'revenue_productivity.png', dpi=300, bbox_inches='tight')
    plt.close()

def plot_cost_analysis(df, output_dir):
    """Plot cost analysis over time."""
    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(12, 10))
    
    if 'TotalCost' in df.columns and 'AvailableBudget' in df.columns:
        # Cost utilization
        ax1.plot(df['TimeStep'], df['TotalCost'], label='Total Cost', 
                linewidth=2, color='red', marker='o', markersize=4)
        ax1.plot(df['TimeStep'], df['AvailableBudget'], label='Available Budget', 
                linewidth=2, color='orange', marker='s', markersize=4)
        ax1.set_xlabel('Time Step')
        ax1.set_ylabel('Cost ($)')
        ax1.set_title('Cost and Budget Over Time')
        ax1.legend()
        ax1.grid(True, alpha=0.3)
        ax1.ticklabel_format(style='plain', axis='y')
        
        # Budget utilization percentage
        budget_utilization = (df['TotalCost'] / (df['TotalCost'] + df['AvailableBudget'])) * 100
        ax2.plot(df['TimeStep'], budget_utilization, label='Budget Utilization (%)', 
                linewidth=2, color='purple', marker='d', markersize=4)
        ax2.set_xlabel('Time Step')
        ax2.set_ylabel('Budget Utilization (%)')
        ax2.set_title('Budget Utilization Over Time')
        ax2.set_ylim(0, 100)
        ax2.grid(True, alpha=0.3)
    
    plt.tight_layout()
    plt.savefig(output_dir / 'cost_analysis.png', dpi=300, bbox_inches='tight')
    plt.close()

def plot_equilibrium_analysis(data, df, output_dir):
    """Plot equilibrium analysis."""
    if data is None:
        print("Warning: Cannot create equilibrium analysis without JSON data")
        return
    
    fig, ((ax1, ax2), (ax3, ax4)) = plt.subplots(2, 2, figsize=(15, 12))
    
    # Time to equilibrium
    if 'TimeToEquilibrium' in data:
        ax1.bar(['Time to Equilibrium'], [data['TimeToEquilibrium']], color='skyblue')
        ax1.set_ylabel('Time Steps')
        ax1.set_title('Time to Reach Equilibrium')
        ax1.grid(True, alpha=0.3)
    
    # Final workforce composition pie chart
    if 'EquilibriumState' in data and 'Workforce' in data['EquilibriumState']:
        workforce = data['EquilibriumState']['Workforce']
        if 'Humans' in workforce and 'AIAgents' in workforce:
            sizes = [workforce['Humans']['Total'], workforce['AIAgents']['Total']]
            labels = ['Human Workers', 'AI Agents']
            colors = ['lightcoral', 'lightskyblue']
            ax2.pie(sizes, labels=labels, colors=colors, autopct='%1.1f%%', startangle=90)
            ax2.set_title('Final Workforce Composition')
    
    # Experience level distribution
    if 'EquilibriumState' in data and 'Workforce' in data['EquilibriumState']:
        workforce = data['EquilibriumState']['Workforce']
        if 'Humans' in workforce and 'ByExperience' in workforce['Humans']:
            exp_data = workforce['Humans']['ByExperience']
            levels = list(exp_data.keys())
            counts = list(exp_data.values())
            ax3.bar(levels, counts, color='lightgreen')
            ax3.set_xlabel('Experience Level')
            ax3.set_ylabel('Number of Workers')
            ax3.set_title('Final Human Experience Distribution')
            ax3.tick_params(axis='x', rotation=45)
            ax3.grid(True, alpha=0.3)
    
    # Catastrophic failures
    if 'TotalCatastrophicFailures' in data:
        ax4.bar(['Total Failures'], [data['TotalCatastrophicFailures']], color='salmon')
        ax4.set_ylabel('Number of Failures')
        ax4.set_title('Total Catastrophic Failures')
        ax4.grid(True, alpha=0.3)
    
    plt.tight_layout()
    plt.savefig(output_dir / 'equilibrium_analysis.png', dpi=300, bbox_inches='tight')
    plt.close()

def create_summary_dashboard(data, df, output_dir):
    """Create a comprehensive dashboard with key metrics."""
    fig = plt.figure(figsize=(20, 12))
    
    # Create a grid layout
    gs = fig.add_gridspec(3, 4, hspace=0.3, wspace=0.3)
    
    # Workforce composition over time
    ax1 = fig.add_subplot(gs[0, :2])
    if 'Workforce.Humans.Total' in df.columns:
        humans_col = 'Workforce.Humans.Total'
        ai_col = 'Workforce.AIAgents.Total'
    else:
        humans_col = 'TotalHumans'
        ai_col = 'TotalAIAgents'
    
    if humans_col in df.columns:
        ax1.plot(df['TimeStep'], df[humans_col], label='Humans', linewidth=2)
        ax1.plot(df['TimeStep'], df[ai_col], label='AI Agents', linewidth=2)
        ax1.set_title('Workforce Evolution')
        ax1.legend()
        ax1.grid(True, alpha=0.3)
    
    # Revenue over time
    ax2 = fig.add_subplot(gs[0, 2:])
    if 'RevenueOutput' in df.columns:
        ax2.plot(df['TimeStep'], df['RevenueOutput'], color='green', linewidth=2)
        ax2.set_title('Revenue Output')
        ax2.grid(True, alpha=0.3)
        ax2.ticklabel_format(style='plain', axis='y')
    
    # Cost utilization
    ax3 = fig.add_subplot(gs[1, :2])
    if 'TotalCost' in df.columns:
        ax3.plot(df['TimeStep'], df['TotalCost'], label='Total Cost', linewidth=2)
        if 'AvailableBudget' in df.columns:
            ax3.plot(df['TimeStep'], df['AvailableBudget'], label='Available Budget', linewidth=2)
        ax3.set_title('Cost Analysis')
        ax3.legend()
        ax3.grid(True, alpha=0.3)
        ax3.ticklabel_format(style='plain', axis='y')
    
    # Productivity
    ax4 = fig.add_subplot(gs[1, 2:])
    if 'TotalProductivity' in df.columns:
        ax4.plot(df['TimeStep'], df['TotalProductivity'], color='blue', linewidth=2)
        ax4.set_title('Total Productivity')
        ax4.grid(True, alpha=0.3)
    
    # Key metrics (bottom row)
    if data:
        # Time to equilibrium
        ax5 = fig.add_subplot(gs[2, 0])
        if 'TimeToEquilibrium' in data:
            ax5.text(0.5, 0.5, f"{data['TimeToEquilibrium']}\nTime Steps", 
                    ha='center', va='center', fontsize=16, weight='bold')
            ax5.set_title('Time to Equilibrium')
            ax5.set_xlim(0, 1)
            ax5.set_ylim(0, 1)
            ax5.axis('off')
        
        # Final workforce ratio
        ax6 = fig.add_subplot(gs[2, 1])
        if 'EquilibriumState' in data and 'Workforce' in data['EquilibriumState']:
            workforce = data['EquilibriumState']['Workforce']
            if 'Humans' in workforce and 'AIAgents' in workforce:
                total_humans = workforce['Humans']['Total']
                total_ai = workforce['AIAgents']['Total']
                ratio = total_ai / total_humans if total_humans > 0 else 0
                ax6.text(0.5, 0.5, f"{ratio:.2f}\nAI:Human Ratio", 
                        ha='center', va='center', fontsize=16, weight='bold')
                ax6.set_title('Final AI:Human Ratio')
                ax6.set_xlim(0, 1)
                ax6.set_ylim(0, 1)
                ax6.axis('off')
        
        # Total failures
        ax7 = fig.add_subplot(gs[2, 2])
        if 'TotalCatastrophicFailures' in data:
            ax7.text(0.5, 0.5, f"{data['TotalCatastrophicFailures']}\nFailures", 
                    ha='center', va='center', fontsize=16, weight='bold')
            ax7.set_title('Catastrophic Failures')
            ax7.set_xlim(0, 1)
            ax7.set_ylim(0, 1)
            ax7.axis('off')
        
        # Final revenue
        ax8 = fig.add_subplot(gs[2, 3])
        if 'EquilibriumState' in data and 'RevenueOutput' in data['EquilibriumState']:
            revenue = data['EquilibriumState']['RevenueOutput']
            ax8.text(0.5, 0.5, f"${revenue:,.0f}\nFinal Revenue", 
                    ha='center', va='center', fontsize=16, weight='bold')
            ax8.set_title('Final Revenue Output')
            ax8.set_xlim(0, 1)
            ax8.set_ylim(0, 1)
            ax8.axis('off')
    
    plt.suptitle('Workforce AI Transition Simulation Dashboard', fontsize=20, weight='bold')
    plt.savefig(output_dir / 'simulation_dashboard.png', dpi=300, bbox_inches='tight')
    plt.close()

def main():
    parser = argparse.ArgumentParser(description='Visualize Workforce AI Transition Simulation results')
    parser.add_argument('input_file', help='Path to simulation results file (JSON or CSV)')
    parser.add_argument('-o', '--output', default='plots', help='Output directory for plots')
    parser.add_argument('--dashboard-only', action='store_true', help='Generate only the dashboard')
    
    args = parser.parse_args()
    
    # Create output directory
    output_dir = Path(args.output)
    output_dir.mkdir(exist_ok=True)
    
    # Load data
    print(f"Loading simulation data from {args.input_file}...")
    data, df = load_simulation_data(args.input_file)
    
    if df is None:
        print("Error: Could not load simulation data")
        sys.exit(1)
    
    print(f"Loaded {len(df)} time steps of simulation data")
    
    if args.dashboard_only:
        print("Creating summary dashboard...")
        create_summary_dashboard(data, df, output_dir)
    else:
        # Generate all plots
        print("Creating workforce composition plots...")
        plot_workforce_composition(df, output_dir)
        
        print("Creating revenue and productivity plots...")
        plot_revenue_and_productivity(df, output_dir)
        
        print("Creating cost analysis plots...")
        plot_cost_analysis(df, output_dir)
        
        print("Creating equilibrium analysis...")
        plot_equilibrium_analysis(data, df, output_dir)
        
        print("Creating summary dashboard...")
        create_summary_dashboard(data, df, output_dir)
    
    print(f"Plots saved to {output_dir}/")
    print("Generated files:")
    for plot_file in output_dir.glob('*.png'):
        print(f"  - {plot_file.name}")

if __name__ == '__main__':
    main()