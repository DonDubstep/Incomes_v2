import Sidebar from './components/Sidebar/Sidebar';
import Header from './components/Header/Header';
import StatsCard from './components/StatsCards/StatsCards';
import CategoriesChart from './components/CategoriesChart/CategoriesChart';
import TransactionsList from './components/TransactionsList/TransactionsList';
import styles from './App.module.css';

function App() {
  return (
    <div className={styles.appContainer}>
      <Sidebar />
      <main className={styles.mainContent}>
        <Header />
        <StatsCard />
        <CategoriesChart />
        <TransactionsList />
      </main>
    </div>
  );
}

export default App;