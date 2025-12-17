package queries

// 1. Get Records. No input
var GetRecords = `SELECT * FROM startup LIMIT 10;`

// 2. Top {10} Funded Industries ever:

var TopIndustries = `
SELECT
    Industry,
    formatReadableQuantity(SUM(Funding_Amount_USD)) AS total_funds
FROM startup
GROUP BY Industry
ORDER BY SUM(Funding_Amount_USD) DESC
LIMIT 10;
`

// 3. Top-funded industries this {year}
// Return the top industries ranked by total startup received in the current year.


var TopIndustries2025 = `
SELECT
    Industry,
    formatReadableQuantity(SUM(Funding_Amount_USD)) AS total_funds
FROM startup
WHERE Year=2025
GROUP BY Industry
ORDER BY SUM(Funding_Amount_USD) DESC;
`

// 4. Most funded startups across all {time}
// Show startups with the highest cumulative funding, regardless of year or city.

var TopFundedStartups = `
SELECT
    Company,
    formatReadableQuantity(SUM(Funding_Amount_USD)) AS total_funds
FROM startup
GROUP BY Company
ORDER BY SUM(Funding_Amount_USD) DESC;
`

// 5. Industry and City
// {time} and {limit} or {city}

var TopCityandIndustries = `
SELECT
    Industry,
    formatReadableQuantity(SUM(Funding_Amount_USD)) AS total_funds,
    City
FROM startup
GROUP BY Industry, City
ORDER BY SUM(Funding_Amount_USD) DESC;
`
