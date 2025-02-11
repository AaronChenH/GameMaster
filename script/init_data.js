/**
 * 数据库初始化脚本 - AaronChenH
 * 版本: 1.1.0
 * 使用说明:
 * 1. 配置数据库连接信息
 * 2. 执行 node init_data.js
 */

// 切换到指定数据库
db = db.getSiblingDB('galaxy_empire_manager');

// 初始化权限数据
const permissions = [
    {
        role: 'admin',
        permissions: ['player_info', 'give_item', 'execute_script', 'manage_users']
    },
    {
        role: 'service',
        permissions: ['player_info', 'give_item']
    }
];

// 插入权限数据
db.permissions.drop();
db.permissions.insertMany(permissions);
print('权限数据初始化完成');

// 初始化管理员账号
const adminUser = {
    username: 'admin',
    // 使用新生成的哈希值
    password: '$2a$10$h67rqyDikF.1f.CliHlHqeMz9vzYzcwUc7p4D6HSLaYVzJ882vt5u',  // admin888
    role: 'admin',
    status: 1,
    created_at: new Date(),
    updated_at: new Date()
};

// 检查是否已存在管理员账号
const existingAdmin = db.users.findOne({username: 'admin'});
if (!existingAdmin) {
    db.users.insertOne(adminUser);
    print('管理员账号初始化完成');
} else {
    print('管理员账号已存在，跳过初始化');
}

// 可以添加其他初始化数据
// 比如：游戏道具列表、VIP等级配置等 