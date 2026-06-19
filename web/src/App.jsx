import AddUser from "./AddUser";
import StatsHistory from "./StatsHistory";

export default function App() {
  return (
    <div className="app">
      <header className="topbar">
        <h1>Habitica Monitor</h1>
        <span className="sub">Мониторинг статов команды</span>
      </header>
      <main className="grid">
        <AddUser />
        <StatsHistory />
      </main>
      <footer className="foot">API: /api/v1</footer>
    </div>
  );
}
