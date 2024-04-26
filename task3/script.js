fetch('https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1')
.then(response => response.json()) //отправляем GET запрос к API и получаем данные в формате JSON
.then(data => {
    const table = document.getElementById('cryptoTable'); 
    data.slice(0, 5).forEach((currency) => { //выбираются первые 5 крипто 
        const row = table.insertRow(); //метод используется для вставки новой строки и дальше устанавливается  id, symbol, name
        row.innerHTML = `
            <td>${currency.id}</td>
            <td>${currency.symbol}</td>
            <td>${currency.name}</td>
        `;
        row.className = 'blue-background';//присваиваем класс, этот класс мы указали внутри style.css
    });

    const usdtSymbol = data.find((currency) => currency.symbol === 'usdt'); //ищем "usdt" крипто
    if (usdtSymbol) {
        const row = table.insertRow(); //если нашли то так же создается строка и дальше устанавливается id, symbol, name
        row.innerHTML = `
            <td>${usdtSymbol.id}</td>
            <td>${usdtSymbol.symbol}</td>
            <td>${usdtSymbol.name}</td>
        `;
        row.className = 'green-background'; //присваиваем класс чтобы сделать зеленый фон
    }
})
.catch(error => console.error('Error when fetching:', error));
