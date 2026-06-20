import { useState } from "react";
import AddUser from "./AddUser";
import StatsHistory from "./StatsHistory";

export default function App() {
  // Счётчик-триггер: растёт при добавлении пользователя,
  // чтобы StatsHistory перезагрузил выпадающий список.
  const [usersVersion, setUsersVersion] = useState(0);

  return (
    <div className="app">
      <header className="topbar">
        <h1>Habitica Monitor</h1>
        <span className="sub">Мониторинг статов команды</span>
      </header>
      <main className="grid">
        <AddUser onCreated={() => setUsersVersion((v) => v + 1)} />
        <StatsHistory usersVersion={usersVersion} />
      </main>
      <footer className="foot">API: /api/v1</footer>
    </div>
  );
}
