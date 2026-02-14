import Plot from 'react-plotly.js';
import { PriceRecord } from '../types';
import { Data, Layout } from 'plotly.js';

interface PriceChartProps {
    data: PriceRecord[];
}

function PriceChart({ data }: PriceChartProps) {
    const actualPrices: Partial<Data> = {
        x: [],
        y: [],
        mode: 'markers',
        name: 'Actual Price',
        marker: { color: '#1f77b4', size: 6 },
    };

    const fairValue: Partial<Data> = {
        x: [],
        y: [],
        mode: 'lines+markers',
        name: 'Fair Value',
        line: { color: '#ff7f0e' },
    };

    data.forEach((d) => {
        if (d.Series === 'daily_price') {
            (actualPrices.x as string[]).push(d.Date);
            (actualPrices.y as number[]).push(d.Price);
        } else if (d.Series === 'fair_value') {
            (fairValue.x as string[]).push(d.Date);
            (fairValue.y as number[]).push(d.Price);
        }
    });

    const layout: Partial<Layout> = {
        xaxis: {
            title: { text: 'Date' }
        },
        yaxis: {
            title: { text: 'Price (USD)' },
            tickprefix: '$'
        },
        margin: { l: 60, r: 30, b: 50, t: 30 },
        legend: { orientation: 'h', y: 1.1 },
        autosize: true,
    };


    return (
        <Plot
            data={[actualPrices as Data, fairValue as Data]}
            layout={layout}
            useResizeHandler={true}
            style={{ width: '100%', height: '75vh' }}
        />
    );
}

export default PriceChart;
