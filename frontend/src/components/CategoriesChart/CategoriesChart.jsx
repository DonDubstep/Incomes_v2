// src/components/CategoriesChart/CategoriesChart.jsx
import styles from './CategoriesChart.module.css';

function CategoriesChart() {
  const categoriesData = [
    { id: 1, name: 'Безналичными', value: '180 000', percentage: 73, color: '#388e3c', icon: '₽' },
    { id: 2, name: 'Наличными',    value: '45 000', percentage: 18, color: '#388e3c', icon: '💵' },
    { id: 3, name: 'В долларах',   value: '15 500', percentage: 6, color: '#388e3c', icon: '💱' },
    { id: 4, name: 'В дирхамах',   value: '4 620', percentage: 3, color: '#388e3c', icon: '💱' },
  ];

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h3 className={styles.title}>Денежные средства</h3>
      </div>
      
      <div className={styles.chartArea}>
        {categoriesData.map((category) => (
          <div key={category.id} className={styles.categoryRow}>
            <div className={styles.categoryInfo}>
              <span className={styles.icon}>{category.icon}</span>
              <span className={styles.name}>{category.name}</span>
            </div>
            
            <div className={styles.valueAndBar}>
              <span className={styles.value}>{category.value}</span>
              <span className={styles.percentage}>{category.percentage}%</span>
              
              <div className={styles.progressBarBg}>
                <div 
                  className={styles.progressBarFill} 
                  style={{ 
                    width: `${category.percentage}%`, 
                    backgroundColor: category.color 
                  }}
                />
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default CategoriesChart;