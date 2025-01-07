import torch
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim
from torch.utils.data import Dataset, DataLoader
import json
import os

# 定义模型
class GobangModel(nn.Module):
    def __init__(self):
        super(GobangModel, self).__init__()

        # 卷积部分
        self.conv1 = nn.Conv2d(1, 16, kernel_size=5, stride=1, padding=0)  # 输入: 1*15*15 -> 输出: 16*11*11
        self.conv2 = nn.Conv2d(16, 32, kernel_size=4, stride=1, padding=1) # 输出: 32*10*10
        self.pool = nn.MaxPool2d(2, 2) # 输出：32*5*5
        self.conv3 = nn.Conv2d(32, 16, kernel_size=3, stride=1, padding=1) # 输出: 16*5*5

        # 全连接部分
        self.fc1 = nn.Linear(16*5*5, 64)
        self.fc2 = nn.Linear(64, 1)

    def forward(self, x):
        # 卷积
        x = self.conv1(x)
        x = F.relu(x)
        x = self.conv2(x)
        x = F.relu(x)
        x = self.pool(x)
        x = self.conv3(x)
        x = F.relu(x)
        x = x.view(-1, 16*5*5)

        # 全连接
        x = self.fc1(x)
        x = F.relu(x)
        x = self.fc2(x)
        x = F.tanh(x)
        return x

# 定义数据集
class GobangDataset(Dataset):
    def __init__(self, min_data, difficulty, noise):
        # AI自我对弈，获取数据集
        r = os.popen(f'generate-data {min_data} {difficulty} {noise}')
        self.data = json.loads(r.read())
        x = [i['Feature'] for i in self.data]
        y = [i['Label'] for i in self.data]
        self.x = torch.tensor(x, dtype=torch.float).view(-1, 1, 15, 15)
        self.y = torch.tensor(y, dtype=torch.float).view(-1, 1)

    def __len__(self):
        return len(self.data)

    def __getitem__(self, idx):
        return self.x[idx], self.y[idx]

if __name__ == '__main__':
    while True:
        # 加载模型
        model = GobangModel()
        if os.path.exists('gobang-model.json'):
            with open('gobang-model.json', 'r') as f:
                loaded_params = json.load(f)
            state_dict = {name: torch.tensor(param) for name, param in loaded_params.items()}
            model.load_state_dict(state_dict)
        else:
            # 保存模型
            model_params = {name: param.detach().numpy().tolist() for name, param in model.state_dict().items()}
            with open("gobang-model.json", "w") as f:
                json.dump(model_params, f)
        # 定义超参数
        epochs = 20
        criterion = nn.MSELoss()
        optimizer = optim.Adam(model.parameters(), lr=0.001)  # Adam优化器
        train_dataset = GobangDataset(min_data=1000, difficulty=1, noise=0.0001)
        train_loader = DataLoader(train_dataset, batch_size=32)
        # 训练模型
        model.train()
        total_loss = 0.0
        for epoch in range(epochs):
            s = 0.0
            for x, y in train_loader:
                # 前向传播
                outputs = model(x)

                # 计算损失
                loss = criterion(outputs, y)

                # 反向传播和优化
                optimizer.zero_grad()  # 清零梯度
                loss.backward()  # 反向传播
                optimizer.step()  # 更新参数

                # 记录损失
                s += loss.item()
            total_loss += s/epochs
        print(total_loss)
        # 保存模型
        model_params = {name: param.detach().numpy().tolist() for name, param in model.state_dict().items()}
        with open("gobang-model.json", "w") as f:
            json.dump(model_params, f)