#!/usr/bin/env python3
"""
Workforce AI Transition Simulator - Interactive Dashboard

This script creates an interactive web-based dashboard using Plotly for exploring
simulation results with interactive controls and real-time updates.

Usage:
    python interactive_dashboard.py simulation_report.json
    python interactive_dashboard.py --sensitivity sensitivity_report.json
"""

import argparse
import json
import pandas as pd
import plotly.graph_objects as go
import plotly.express as px
from plotly.subplots import make_subplots
import plotly.offline as pyo
from pathlib import Path
import sys

def load_simulation_data(file_path):
    """Load simulation data from JSON file."""
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
    else:
        print(f"Error: Unsupported file format {file_path.suffix}")
        return None, None

def create_simulation_dashboard(data, df, output_path):
    """Create interactive dashboard for simulation results."""
    
    # Create subplots
    fig = make_subplots(
        rows=3, cols=2,
        subplot_titles=('Workforce Composition', 'Revenue Output', 
                       'Cost Analysis', 'Productivity', 
                       'Budget Utilization', 'Key Metrics'),
        specs=[[{"secondary_y": False}, {"secondary_y": False}],
               [{"secondary_y": True}, {"secondary_y": False}],
               [{"secondary_y": False}, {"type": "table"}]]
    )
    
    # Workforce composition
    if 'Workforce.Humans.Total' in df.columns:
        humans_col = 'Workforce.Humans.Total'
        ai_col = 'Workforce.AIAgents.Total'
    else:
        humans_col = 'TotalHumans'
        ai_col = 'TotalAIAgents'
    
    if humans_col in df.columns:
        fig.add_trace(
            go.Scatter(x=df['TimeStep'], y=df[humans_col], 
                      name='Human Workers', line=dict(color='blue', width=3),
                      hovertemplate='Time Step: %{x}<br>Humans: %{y}<extra></extra>'),
            row=1, col=1
        )
        fig.add_trace(
            go.Scatter(x=df['TimeStep'], y=df[ai_col], 
                      name='AI Agents', line=dict(color='red', width=3),
                      hovertemplate='Time Step: %{x}<br>AI Agents: %{y}<extra></extra>'),
            row=1, col=1
        )
    
    # Revenue output
    if 'RevenueOutput' in df.columns:
        fig.add_trace(
            go.Scatter(x=df['TimeStep'], y=df['RevenueOutput'], 
                      name='Revenue', line=dict(color='green', width=3),
                      hovertemplate='Time Step: %{x}<br>Revenue: $%{y:,.0f}<extra></extra>'),
            row=1, col=2
        )
    
    # Cost analysis
    if 'TotalCost' in df.columns:
        fig.add_trace(
            go.Scatter(x=df['TimeStep'], y=df['TotalCost'], 
                      name='Total Cost', line=dict(color='orange', width=3),
                      hovertemplate='Time Step: %{x}<br>Cost: $%{y:,.0f}<extra></extra>'),
            row=2, col=1
        )
    
    if 'AvailableBudget' in df.columns:
        fig.add_trace(
            go.Scatter(x=df['TimeStep'], y=df['AvailableBudget'], 
                      name='Available Budget', line=dict(color='purple', width=3, dash='dash'),
                      hovertemplate='Time Step: %{x}<br>Available: $%{y:,.0f}<extra></extra>'),
            row=2, col=1, secondary_y=True
        )
    
    # Productivity
    if 'TotalProductivity' in df.columns:
        fig.add_trace(
            go.Scatter(x=df['TimeStep'], y=df['TotalProductivity'], 
                      name='Productivity', line=dict(color='teal', width=3),
                      hovertemplate='Time Step: %{x}<br>Productivity: %{y:.2f}<extra></extra>'),
            row=2, col=2
        )
    
    # Budget utilization
    if 'TotalCost' in df.columns and 'AvailableBudget' in df.columns:
        budget_util = (df['TotalCost'] / (df['TotalCost'] + df['AvailableBudget'])) * 100
        fig.add_trace(
            go.Scatter(x=df['TimeStep'], y=budget_util, 
                      name='Budget Utilization', line=dict(color='crimson', width=3),
                      hovertemplate='Time Step: %{x}<br>Utilization: %{y:.1f}%<extra></extra>'),
            row=3, col=1
        )
    
    # Key metrics table
    if data:
        metrics_data = []
        if 'TimeToEquilibrium' in data:
            metrics_data.append(['Time to Equilibrium', f"{data['TimeToEquilibrium']} steps"])
        if 'TotalCatastrophicFailures' in data:
            metrics_data.append(['Catastrophic Failures', str(data['TotalCatastrophicFailures'])])
        if 'EquilibriumState' in data:
            eq_state = data['EquilibriumState']
            if 'RevenueOutput' in eq_state:
                metrics_data.append(['Final Revenue', f"${eq_state['RevenueOutput']:,.0f}"])
            if 'TotalProductivity' in eq_state:
                metrics_data.append(['Final Productivity', f"{eq_state['TotalProductivity']:.2f}"])
        
        if metrics_data:
            fig.add_trace(
                go.Table(
                    header=dict(values=['Metric', 'Value'], 
                               fill_color='lightblue',
                               font=dict(size=14, color='black')),
                    cells=dict(values=list(zip(*metrics_data)),
                              fill_color='white',
                              font=dict(size=12))
                ),
                row=3, col=2
            )
    
    # Update layout
    fig.update_layout(
        title=dict(
            text='Workforce AI Transition Simulation Dashboard',
            x=0.5,
            font=dict(size=24, color='darkblue')
        ),
        height=1000,
        showlegend=True,
        legend=dict(x=0.01, y=0.99),
        template='plotly_white'
    )
    
    # Update axes labels
    fig.update_xaxes(title_text="Time Step", row=1, col=1)
    fig.update_xaxes(title_text="Time Step", row=1, col=2)
    fig.update_xaxes(title_text="Time Step", row=2, col=1)
    fig.update_xaxes(title_text="Time Step", row=2, col=2)
    fig.update_xaxes(title_text="Time Step", row=3, col=1)
    
    fig.update_yaxes(title_text="Number of Workers", row=1, col=1)
    fig.update_yaxes(title_text="Revenue ($)", row=1, col=2)
    fig.update_yaxes(title_text="Cost ($)", row=2, col=1)
    fig.update_yaxes(title_text="Productivity", row=2, col=2)
    fig.update_yaxes(title_text="Budget Utilization (%)", row=3, col=1)
    
    # Save as HTML
    pyo.plot(fig, filename=str(output_path), auto_open=False)
    print(f"Interactive dashboard saved to {output_path}")

def create_sensitivity_dashboard(data, output_path):
    """Create interactive dashboard for sensitivity analysis results."""
    
    if not data or 'ParameterRankings' not in data:
        print("Error: No parameter rankings found in sensitivity data")
        return
    
    rankings = data['ParameterRankings']
    
    # Create subplots
    fig = make_subplots(
        rows=2, cols=2,
        subplot_titles=('Time to Equilibrium Impact', 'Workforce Composition Impact',
                       'Parameter Variations', 'Summary Statistics'),
        specs=[[{"type": "bar"}, {"type": "bar"}],
               [{"type": "scatter"}, {"type": "table"}]]
    )
    
    # Parameter impact on time to equilibrium
    time_impacts = []
    comp_impacts = []
    param_names = []
    
    for param in rankings:
        param_names.append(param['ParameterName'])
        time_impacts.append(param.get('TimeToEquilibriumImpact', 0))
        comp_impacts.append(param.get('WorkforceCompositionImpact', 0))
    
    fig.add_trace(
        go.Bar(x=param_names, y=time_impacts, name='Time Impact',
               marker_color='skyblue',
               hovertemplate='Parameter: %{x}<br>Impact: %{y:.3f}<extra></extra>'),
        row=1, col=1
    )
    
    fig.add_trace(
        go.Bar(x=param_names, y=comp_impacts, name='Composition Impact',
               marker_color='lightcoral',
               hovertemplate='Parameter: %{x}<br>Impact: %{y:.3f}<extra></extra>'),
        row=1, col=2
    )
    
    # Parameter variations
    if 'SensitivityResults' in data:
        results = data['SensitivityResults']
        colors = px.colors.qualitative.Set1
        
        for i, param_result in enumerate(results):
            param_name = param_result['ParameterName']
            color = colors[i % len(colors)]
            
            if 'ParameterValues' in param_result and 'TimeToEquilibriumByValue' in param_result:
                param_values = param_result['ParameterValues']
                time_values = [param_result['TimeToEquilibriumByValue'].get(str(v), 0) 
                              for v in param_values]
                
                fig.add_trace(
                    go.Scatter(x=param_values, y=time_values, 
                              name=param_name, mode='lines+markers',
                              line=dict(color=color, width=3),
                              marker=dict(size=8),
                              hovertemplate=f'{param_name}: %{{x}}<br>Time to Equilibrium: %{{y}}<extra></extra>'),
                    row=2, col=1
                )
    
    # Summary statistics table
    summary_data = [
        ['Total Parameters Analyzed', str(len(param_names))],
        ['Most Impactful (Time)', param_names[time_impacts.index(max(time_impacts))] if time_impacts else 'N/A'],
        ['Most Impactful (Composition)', param_names[comp_impacts.index(max(comp_impacts))] if comp_impacts else 'N/A'],
        ['Avg Time Impact', f"{sum(time_impacts)/len(time_impacts):.3f}" if time_impacts else 'N/A'],
        ['Avg Composition Impact', f"{sum(comp_impacts)/len(comp_impacts):.3f}" if comp_impacts else 'N/A']
    ]
    
    fig.add_trace(
        go.Table(
            header=dict(values=['Metric', 'Value'], 
                       fill_color='lightgreen',
                       font=dict(size=14, color='black')),
            cells=dict(values=list(zip(*summary_data)),
                      fill_color='white',
                      font=dict(size=12))
        ),
        row=2, col=2
    )
    
    # Update layout
    fig.update_layout(
        title=dict(
            text='Sensitivity Analysis Dashboard',
            x=0.5,
            font=dict(size=24, color='darkgreen')
        ),
        height=800,
        showlegend=True,
        template='plotly_white'
    )
    
    # Update axes
    fig.update_xaxes(title_text="Parameter", row=1, col=1)
    fig.update_xaxes(title_text="Parameter", row=1, col=2)
    fig.update_xaxes(title_text="Parameter Value", row=2, col=1)
    
    fig.update_yaxes(title_text="Impact Score", row=1, col=1)
    fig.update_yaxes(title_text="Impact Score", row=1, col=2)
    fig.update_yaxes(title_text="Time to Equilibrium", row=2, col=1)
    
    # Save as HTML
    pyo.plot(fig, filename=str(output_path), auto_open=False)
    print(f"Interactive sensitivity dashboard saved to {output_path}")

def main():
    parser = argparse.ArgumentParser(description='Create Interactive Workforce AI Transition Dashboard')
    parser.add_argument('input_file', help='Path to simulation results file (JSON)')
    parser.add_argument('-o', '--output', help='Output HTML file path')
    parser.add_argument('--sensitivity', action='store_true', help='Create sensitivity analysis dashboard')
    parser.add_argument('--open', action='store_true', help='Open dashboard in browser after creation')
    
    args = parser.parse_args()
    
    # Determine output path
    if args.output:
        output_path = Path(args.output)
    else:
        input_path = Path(args.input_file)
        if args.sensitivity:
            output_path = input_path.parent / f"{input_path.stem}_sensitivity_dashboard.html"
        else:
            output_path = input_path.parent / f"{input_path.stem}_dashboard.html"
    
    # Load data
    print(f"Loading data from {args.input_file}...")
    
    if args.sensitivity:
        # Load sensitivity analysis data
        with open(args.input_file, 'r') as f:
            data = json.load(f)
        create_sensitivity_dashboard(data, output_path)
    else:
        # Load simulation data
        data, df = load_simulation_data(args.input_file)
        if data is None or df is None:
            print("Error: Could not load simulation data")
            sys.exit(1)
        create_simulation_dashboard(data, df, output_path)
    
    if args.open:
        import webbrowser
        webbrowser.open(f'file://{output_path.absolute()}')
        print(f"Dashboard opened in browser")

if __name__ == '__main__':
    main()