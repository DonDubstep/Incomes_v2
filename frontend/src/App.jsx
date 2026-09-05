import Sidebar from './components/Sidebar/Sidebar';
import Header from './components/Header/Header';
import StatsCard from './components/StatsCards/StatsCards';
import CategoriesChart from './components/CategoriesChart/CategoriesChart';
import TransactionsList from './components/TransactionsList/TransactionsList';
import styles from './App.module.css';

function App() {
  // Данные для карточек статистики (как будто с бэкенда)
  const statsData = [
    { 
      id: 1, 
      title: 'Общий доход', 
      value: '245 120 ₽', 
      trend: '+12.5% по сравнению с апрелем', 
      isLarge: true 
    },
    { 
      id: 2, 
      title: 'Средний доход в день', 
      value: '7 907 ₽' 
    },
    { 
      id: 3, 
      title: 'Количество транзакций', 
      value: '28' 
    },
  ];

  return (
    <div className={styles.appContainer}>
      <Sidebar />
      <main className={styles.mainContent}>
        <Header />
        
        {/* Рендерим карточки статистики через map() */}
        <div className={styles.statsGrid}>
          {statsData.map((stat) => (
            <StatsCard 
              key={stat.id} 
              title={stat.title} 
              value={stat.value} 
              trend={stat.trend} 
              isLarge={stat.isLarge} 
            />
          ))}
        </div>

        <CategoriesChart />
        <TransactionsList />
      </main>
    </div>
  );
}

export default App;