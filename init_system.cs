using System;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;

class Program
{
    static async Task Main(string[] args)
    {
        var client = new HttpClient();
        var baseURL = "http://localhost:8890/api/v1/public/init";

        var payload = new
        {
            admin = new
            {
                username = "admin",
                password = "admin123",
                email = "admin@example.com"
            },
            user = new
            {
                username = "user",
                password = "user123",
                email = "user@example.com"
            },
            database = new
            {
                type = "mysql",
                host = "localhost",
                port = "3306",
                database = "oneclickvirt",
                username = "oneclickvirt",
                password = "123456"
            }
        };

        var json = System.Text.Json.JsonSerializer.Serialize(payload);
        var content = new StringContent(json, Encoding.UTF8, "application/json");

        try
        {
            Console.WriteLine("开始系统初始化...");
            var response = await client.PostAsync(baseURL, content);
            var responseContent = await response.Content.ReadAsStringAsync();

            if (response.IsSuccessStatusCode)
            {
                Console.WriteLine("✅ 系统初始化成功！");
                Console.WriteLine($"Response: {responseContent}");
            }
            else
            {
                Console.WriteLine($"❌ 初始化失败: {response.StatusCode}");
                Console.WriteLine($"Response: {responseContent}");
            }
        }
        catch (Exception ex)
        {
            Console.WriteLine($"❌ 发生错误: {ex.Message}");
        }
    }
}
