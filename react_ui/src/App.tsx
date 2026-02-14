import { useState, useEffect } from 'react';
import Header from './components/Header';
import Sidebar from './components/Sidebar';
import PriceChart from './components/PriceChart'
import { PriceRecord } from './types';
import './App.css';

function App() {
    const [data, setData] = useState<PriceRecord[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [isSidebarOpen, setIsSidebarOpen] = useState(true);

    useEffect(() => {
        // TODO pull this out into its own .ts file?
        fetch('/api/data')
            .then((res) => {
                if (!res.ok) {
                    throw new Error(`HTTP error! status: ${res.status}`);
                }
                return res.json();
            })
            .then((fetchedData: PriceRecord[]) => {
                setData(fetchedData);
                setIsLoading(false);
            })
            .catch((err) => {
                setError(err.message);
                setIsLoading(false);
            });
    }, []);

    const ticker = data.length > 0 ? data[0].Ticker : 'N/A';

    return (
        <>
            <Header />
            <div className="app-body">
                <Sidebar isOpen={isSidebarOpen} onToggle={() => setIsSidebarOpen(!isSidebarOpen)} />
                <main className={`main-content ${isSidebarOpen ? '' : 'sidebar-collapsed'}`}>
                    {isLoading && <p>Loading chart data...</p>}
                    {error && <p className="error-message">Error: {error}</p>}
                    {!isLoading && !error && (
                        <>
                            <div className="chart-header">
                                <h2>Lynch Fair Value Analysis for {ticker}</h2>
                            </div>
                            <PriceChart data={data} />
                        </>
                    )}
                </main>
            </div>
        </>
    );
}

export default App;
