CREATE TABLE Users (
                       UserID INT IDENTITY(1,1) PRIMARY KEY,
                       Username NVARCHAR(50) NOT NULL UNIQUE,
                       Email NVARCHAR(100) NOT NULL UNIQUE,
                       PasswordHash NVARCHAR(255) NOT NULL,
                       Role NVARCHAR(50) NULL, -- e.g. admin, staff, manager
                       IsActive BIT NOT NULL DEFAULT 1,
                       EmployeeID INT NULL,

                       CreatedAt DATETIME NOT NULL DEFAULT GETDATE(),
                       UpdatedAt DATETIME NOT NULL DEFAULT GETDATE(),

                       CONSTRAINT FK_Users_Employees FOREIGN KEY (EmployeeID)
                           REFERENCES Employees(EmployeeID)
                           ON DELETE SET NULL
)
    GO

-- Index untuk pencarian cepat
CREATE INDEX IX_Users_Username ON Users(Username);
CREATE INDEX IX_Users_Email ON Users(Email);