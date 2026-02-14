import './Sidebar.css';

interface SidebarProps {
    isOpen: boolean;
    onToggle: () => void;
}

function Sidebar({ isOpen, onToggle }: SidebarProps) {
    return (
        <aside className={`sidebar ${isOpen ? 'open' : 'collapsed'}`}>
            <button className="sidebar-toggle" onClick={onToggle} aria-label="Toggle Sidebar">
                {isOpen ? '◀' : '▶'}
            </button>
            {isOpen && (
                <nav className="sidebar-nav">
                    <h3>Data Sources</h3>
                    <p className="placeholder-text">
                        UNDER CONSTRUCTION: This panel will be used to browse Parquet files and select datasets.
                    </p>
                    <ul>
                        <li>PLACEHOLDER: AAPL.parquet</li>
                    </ul>
                </nav>
            )}
        </aside>
    );
}

export default Sidebar;
