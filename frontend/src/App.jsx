import { useEffect, useState } from 'react';
import Sidebar from './components/Sidebar/Sidebar';
import Header from './components/Header/Header';
import StatsCard from './components/StatsCards/StatsCards';
import CategoriesChart from './components/CategoriesChart/CategoriesChart';
import TransactionsList from './components/TransactionsList/TransactionsList';
import StatisticsPage from './components/StatisticsPage/StatisticsPage';
import styles from './App.module.css';

// Простая навигация через hash: #/stats - статистика, иначе - главная
const pageFromHash = () => (window.location.hash === '#/stats' ? 'stats' : 'incomes');

function App() {
  const [page, setPage] = useState(pageFromHash);

  useEffect(() => {
    const onHashChange = () => setPage(pageFromHash());
    window.addEventListener('hashchange', onHashChange);
    return () => window.removeEventListener('hashchange', onHashChange);
  }, []);

  const navigate = (next) => {
    window.location.hash = next === 'stats' ? '#/stats' : '#/';
  };

  return (
    <div className={styles.appContainer}>
      <Sidebar activePage={page} onNavigate={navigate} />
      <main className={styles.mainContent}>
        {page === 'stats' ? (
          <StatisticsPage />
        ) : (
          <>
            <Header />
            <StatsCard />
            <CategoriesChart />
            <TransactionsList />
          </>
        )}
      </main>
    </div>
  );
}

export default App;
