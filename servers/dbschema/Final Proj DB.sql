-- Note: since MSSQL is large
-- We do not dockerize it into EC2 (1 GB space is not enough)
-- Instead MSSQL is hosted in remote amazon MSSQL services
-- this schema shows the RoomReservation in the remote database

CREATE DATABASE RoomReservation

GO

USE RoomReservation

GO

CREATE TABLE tblUserType (
    userTypeID INT PRIMARY KEY IDENTITY(1, 1) NOT NULL,
    userTypeName VARCHAR(32) NOT NULL
)

GO

CREATE FUNCTION fn_noDuplicateUserType()
RETURNS INT
AS
BEGIN
DECLARE @Return INT = 0
IF EXISTS (SELECT userTypeName FROM tblUserType GROUP BY userTypeName HAVING COUNT(userTypeName) > 1)
    SET @Return = 1
RETURN @Return
END

GO


-- Constraint for unique user type
ALTER TABLE tblUserType
ADD CONSTRAINT noDuplicateUserType
CHECK([dbo].fn_noDuplicateUserType() = 0)

GO


INSERT INTO tblUserType(userTypeName)
VALUES('Admin'), ('Normal')

GO

CREATE TABLE tblUser (
    userID INT PRIMARY KEY IDENTITY(1, 1) NOT NULL,
    userName VARCHAR(64) NOT NULL,
    email VARCHAR(128) NOT NULL,
    passHash BINARY(60) NOT NULL,
    userTypeID INT FOREIGN KEY REFERENCES tblUserType(userTypeID) NOT NULL
)

GO



-- Business Rule: user name in tblUser must be unique
-- return 1 (Error) if the business rule is violated
-- return 0 if not
CREATE FUNCTION fn_noDuplicateUserName()
RETURNS INT
AS
BEGIN
DECLARE @Return INT = 0
IF EXISTS(SELECT userName
    FROM tblUser 
    GROUP BY userName
    HAVING COUNT(userName) > 1)
    SET @Return = 1

RETURN @Return
END
GO


-- Business Rule: email in tblUser must be unique for each user
-- return 1 (Error) if the business rule is violated
-- return 0 if not
CREATE FUNCTION fn_noDuplicateEmail()
RETURNS INT
AS
BEGIN
DECLARE @Return INT = 0
IF EXISTS(SELECT email
    FROM tblUser 
    GROUP BY email
    HAVING COUNT(email) > 1)
    SET @Return = 1

RETURN @Return
END

GO

-- add the business rule of no duplicate username in tblUser
ALTER TABLE tblUser
ADD CONSTRAINT noDuplicateUserName
CHECK([dbo].fn_noDuplicateUserName() = 0)

GO

-- add the business rule of no duplicate email in tblUser
ALTER TABLE tblUser
ADD CONSTRAINT noDuplicateEMail
CHECK([dbo].fn_noDuplicateEmail() = 0)

GO


CREATE TABLE tblRoomType (
    roomTypeID INT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
    roomTypeName VARCHAR(64) NOT NULL
)

GO

CREATE FUNCTION fn_noDuplicateRoomType()
RETURNS INT
AS
BEGIN
DECLARE @Return INT = 0
IF EXISTS (SELECT roomTypeName FROM tblRoomType GROUP BY roomTypeName HAVING COUNT(roomTypeName) > 1)
    SET @Return = 1
RETURN @Return
END


GO

-- Constraint for unique room type
ALTER TABLE tblRoomType
ADD CONSTRAINT noDuplicateRoomType
CHECK([dbo].fn_noDuplicateRoomType() = 0)

GO


INSERT INTO tblRoomType(roomTypeName)
VALUES('Study'), ('Teamwork'), ('Demonstration'), ('Lounge'), ('Computer Lab'), ('Other')

GO

CREATE TABLE tblRoomStatus(
    roomStatusID INT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
    roomStatusName VARCHAR(64) NOT NULL
)

GO

CREATE FUNCTION fn_noDuplicateRoomStatus()
RETURNS INT
AS
BEGIN
DECLARE @Return INT = 0
IF EXISTS (SELECT roomStatusName FROM tblRoomStatus GROUP BY roomStatusName HAVING COUNT(roomStatusName) > 1)
    SET @Return = 1
RETURN @Return
END

GO

-- Constraint for unique room type
ALTER TABLE tblRoomStatus
ADD CONSTRAINT noDuplicateRoomStatus
CHECK([dbo].fn_noDuplicateRoomStatus() = 0)

GO

INSERT INTO tblRoomStatus(roomStatusName)
VALUES('Functioning'), ('Malfunctioning'), ('Cleaning'), ('Maintaining'), ('Building'), ('Renovating')

GO

CREATE TABLE tblRoom(
    roomID INT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
    roomName VARCHAR(64) NOT NULL,
    roomFloor INT NOT NULL,
    capacity INT NOT NULL,
    roomTypeID INT FOREIGN KEY REFERENCES tblRoomType(roomTypeID) NOT NULL
)

GO

-- constraint 1: capacity must be positive
-- constraint 2: for teamwork/demonstration room, capacity must be at least 2
-- constraint 3: unique room name

CREATE FUNCTION fn_capacityMustBePositive()
RETURNS INT
AS
BEGIN
DECLARE @Return INT = 0
IF EXISTS(select * from tblRoom WHERE capacity <= 0)
    SET @Return = 1

RETURN @Return
END

GO

ALTER TABLE tblRoom
ADD CONSTRAINT capacityMustBePositive
CHECK([dbo].fn_capacityMustBePositive() = 0)

GO

CREATE FUNCTION fn_teamRoomConstraint()
RETURNS INT
AS
BEGIN
DECLARE @Return INT = 0
IF EXISTS (
    SElECT *
    FROM tblRoom R JOIN tblRoomType RT ON R.roomTypeID = RT.roomTypeID
    WHERE RT.roomTypeName = 'Teamwork' AND R.capacity <= 1
)
SET @Return = 1
RETURN @Return
END

GO

ALTER TABLE tblRoom
ADD CONSTRAINT teamRoomConstraint 
CHECK([dbo].fn_teamRoomConstraint() = 0)

GO

CREATE FUNCTION fn_demoRoomConstraint()
RETURNS INT
AS
BEGIN
DECLARE @Return INT = 0
IF EXISTS (
    SElECT *
    FROM tblRoom R JOIN tblRoomType RT ON R.roomTypeID = RT.roomTypeID
    WHERE RT.roomTypeName = 'Demonstration' AND R.capacity <= 1
)
SET @Return = 1
RETURN @Return
END

GO

ALTER TABLE tblRoom
ADD CONSTRAINT demoRoomConstraint 
CHECK([dbo].fn_demoRoomConstraint() = 0)

GO

CREATE FUNCTION fn_uniqueRoomName()
RETURNS INT
AS
BEGIN
DECLARE @Return INT = 0
IF EXISTS(SELECT roomName FROM tblRoom GROUP BY roomName HAVING COUNT(roomName) > 1)
SET @Return = 1
RETURN @Return
END

GO

ALTER TABLE tblRoom
ADD CONSTRAINT uniqueRoomNameConstraint
CHECK([dbo].fn_uniqueRoomName() = 0)
GO

CREATE TABLE tblRoomStatusLog(
    roomStatusLogID INT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
    roomID INT FOREIGN KEY REFERENCES tblRoom(roomID) NOT NULL,
    roomStatusID INT FOREIGN KEY REFERENCES tblRoomStatus(roomStatusID) NOT NULL,
    beginDate DATE NOT NULL,
    endDate DATE
)

GO

CREATE TABLE tblReservation(
    reservationID INT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
    userID INT FOREIGN KEY REFERENCES tblUser(userID) NOT NULL,
    roomID INT FOREIGN KEY REFERENCES tblRoom(roomID) NOT NULL,
    tranDate DATE NOT NULL,
    reserveDate DATE NOT NULL,
    beginTime INT NOT NULL,
    duration INT NOT NULL
)

GO

ALTER TABLE tblReservation
ADD endTime AS (beginTime + duration)

GO

-- 1:  no res before 8.

CREATE FUNCTION fn_NoStartTimeBefore8()
RETURNS INT
AS
BEGIN
DECLARE @Result INT = 0
IF EXISTS(SELECT * FROM tblReservation WHERE beginTime < 16)
    SET @Result = 1
RETURN @Result
END

GO

ALTER TABLE tblReservation
ADD CONSTRAINT noStartTimeBefore8 
CHECK([dbo].fn_NoStartTimeBefore8() = 0)

GO

-- 2: no more res after 9

CREATE FUNCTION fn_NoTimeAfter9()
RETURNS INT
AS
BEGIN
DECLARE @Result INT = 0
IF EXISTS(SELECT * FROM tblReservation WHERE beginTime + duration > 42)
    SET @Result = 1
RETURN @Result
END

GO

ALTER TABLE tblReservation
ADD CONSTRAINT noTimeAfter9 
CHECK([dbo].fn_NoTimeAfter9() = 0)

GO

-- 3: reserve time must be positive

CREATE FUNCTION fn_PositiveResTime()
RETURNS INT
AS
BEGIN 
DECLARE @Result INT = 0
IF EXISTS(SELECT * FROM tblReservation WHERE duration <= 0)
    SET @Result = 1
RETURN @Result
END

GO

ALTER TABLE tblReservation
ADD CONSTRAINT PositiveResTime 
CHECK([dbo].fn_PositiveResTime() = 0)


GO

-- 4: reserve time must be at least current
CREATE FUNCTION fn_ResAtCurrent()
RETURNS INT
AS
BEGIN
DECLARE @Result INT = 0
IF EXISTS(SELECT * FROM tblReservation WHERE reserveDate < GETDATE())
    SET @Result = 1
RETURN @Result
END

GO

ALTER TABLE tblReservation
ADD CONSTRAINT ResAtCurrent
CHECK([dbo].fn_ResAtCurrent() = 0)

GO
-- 5: user can reserve at most 5 hours
CREATE FUNCTION fn_NoMore5HoursForUser()
RETURNS INT
AS
BEGIN
DECLARE @Return INT = 0
IF EXISTS(SELECT userID
            FROM tblReservation
            GROUP BY userID, reserveDate
            HAVING SUM(duration) > 10)
    SET @Return = 1
RETURN @Return
END

GO

ALTER TABLE tblReservation
ADD CONSTRAINT noMore5HrForAUser
CHECK([dbo].fn_NoMore5HoursForUser() = 0)

GO


-- 6: Admin cannot reserve room?


GO

CREATE TABLE tblEquipment (
    equipID INT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
    equipName VARCHAR(128) NOT NULL
)

GO

-- Unique EquipmentName
CREATE FUNCTION fn_noDuplicateEquip()
RETURNS INT
AS
BEGIN
DECLARE @Return INT = 0
IF EXISTS (SELECT equipName FROM tblEquipment GROUP BY equipName HAVING COUNT(equipName) > 1)
    SET @Return = 1
RETURN @Return
END

GO

ALTER TABLE tblEquipment
ADD CONSTRAINT noDuplicateEquip
CHECK([dbo].fn_noDuplicateEquip() = 0)

GO

CREATE TABLE tblEquipInRoom (
    equipRoomID INT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
    equipID INT FOREIGN KEY REFERENCES tblEquipment(equipID) NOT NULL,
    roomID INT FOREIGN KEY REFERENCES tblRoom(roomID) NOT NULL,
    addDate DATE NOT NULL,
    removeDate DATE
)

GO


CREATE TABLE tblRoomIssue(
    roomIssueID INT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
    roomID INT FOREIGN KEY REFERENCES tblRoom(roomID) NOT NULL,
    roomIssueBody VARCHAR(256) NOT NULL,
    createDate DATE NOT NULL,
    confirmDate DATE,
    solveDate DATE
)

GO
CREATE FUNCTION fn_noSolveBeforeConfirm()
RETURNS INT
AS
BEGIN
DECLARE @Return INT = 0
IF EXISTS(SELECT * FROM tblRoomIssue WHERE confirmDate IS NULL AND solveDate IS NOT NULL)
    SET @Return = 1
RETURN @Return
END

GO

ALTER TABLE tblRoomIssue
ADD CONSTRAINT noSolveBeforeConform
CHECK([dbo].fn_noSolveBeforeConfirm() = 0)

GO


CREATE PROCEDURE usp_getUserTypeID
@UTypeName VARCHAR(64),
@UTypeID INT OUTPUT
AS
SET @UTypeID = (SELECT TOP 1 userTypeID FROM tblUserType WHERE userTypeName = @UTypeName)

GO

CREATE PROCEDURE usp_createUser
@userName VARCHAR(64),
@email VARCHAR(128),
@passHash BINARY(60),
@userTypeName VARCHAR(64)
AS
IF @userName IS NULL OR @email IS NULL OR @passHash IS NULL OR @userTypeName IS NULL
    THROW 50001, 'Some Params is NULL', 1
DECLARE @userTypeID INT
EXEC usp_getUserTypeID
@UTypeName = @userTypeName,
@UTypeID = @userTypeID OUTPUT

IF @userTypeID IS NULL
    THROW 50002, 'User Type Not Found', 1
BEGIN TRAN
INSERT INTO tblUser(userName, email, passHash, userTypeID)
VALUES(@userName, @email, @passHash, @userTypeID)
IF @@ERROR<>0
    ROLLBACK TRAN
ELSE
    COMMIT TRAN

GO

CREATE PROCEDURE usp_getRoomTypeID
@RTypeName VARCHAR(64),
@RTypeID INT OUTPUT
AS
SET @RTypeID = (SELECT TOP 1 roomTypeID FROM tblRoomType WHERE roomTypeName = @RTypeName)

GO

CREATE PROCEDURE usp_getUserUserType
@UName VARCHAR(64),
@UTName VARCHAR(64) OUTPUT
AS
SET @UTName = (SELECT TOP 1 UT.userTypeName 
                    FROM tblUserType UT JOIN tblUser U ON U.userTypeID = UT.userTypeID 
                    WHERE U.userName = @UName)



GO

CREATE PROCEDURE usp_createRoom
@roomName VARCHAR(64),
@floor INT,
@capcity INT,
@roomTypeName VARCHAR(64),
@userName VARCHAR(64)
AS
IF @roomName IS NULL OR @floor IS NULL OR @capcity IS NULL OR @roomTypeName IS NULL OR @userName IS NULL
    THROW 50001, 'Some Params are NULL', 1
DECLARE @RoomTypeID INT
EXEC usp_getRoomTypeID
@RTypeName = @roomTypeName,
@RTypeID = @RoomTypeID OUTPUT

IF @RoomTypeID IS NULL
    THROW 50002, 'Room Type NOT Found', 1

DECLARE @curUserType VARCHAR(64)
EXEC usp_getUserUserType
@UName = @userName,
@UTName = @curUserType OUTPUT

IF @curUserType != 'Admin'
    THROW 50002, 'This user is not admin', 1

BEGIN TRAN
INSERT INTO tblRoom(roomName, roomFloor, capacity, roomTypeID)
VALUES(@roomName, @floor, @capcity, @RoomTypeID)
IF @@ERROR<>0
    ROLLBACK TRAN
ELSE
    COMMIT TRAN

GO

CREATE PROCEDURE usp_getRoomID
@RName VARCHAR(64),
@RID INT OUTPUT
AS
SET @RID = (SELECT TOP 1 roomID FROM tblRoom WHERE roomName = @RName)


GO

CREATE PROCEDURE usp_getUserID
@UName VARCHAR(64),
@UID INT OUTPUT
AS
SET @UID = (SELECT TOP 1 userID FROM tblUser WHERE userName = @UName)

GO

CREATE PROCEDURE usp_makeRoomReservation
@userName VARCHAR(64),
@roomName VARCHAR(64),
@tranDate DATE,
@reserveDate DATE,
@beginTime INT,
@duration INT
AS
IF @userName IS NULL OR @roomName IS NULL OR @tranDate IS NULL OR @reserveDate IS NULL OR @beginTime IS NULL OR @duration IS NULL
    THROW 50001, 'Some Params are null', 1
DECLARE @userID INT
EXEC usp_getUserID
@UName = @userName,
@UID = @userID OUTPUT

IF @userID IS NULL 
    THROW 50002, 'User Not Found', 1

DECLARE @RoomID INT
EXEC usp_getRoomID
@RName = @roomName,
@RID = @RoomID OUTPUT

IF @RoomID IS NULL
    THROW 50002, 'Room Not Found', 1

BEGIN TRAN
INSERT INTO tblReservation(userID, roomID, tranDate, reserveDate, beginTime, duration)
VALUES(@userID, @RoomID, @tranDate, @reserveDate, @beginTime, @duration)
IF @@ERROR<>0
    ROLLBACK TRAN
ELSE
    COMMIT TRAN

GO

CREATE PROCEDURE usp_addEquipment
@equipName VARCHAR(128),
@userName VARCHAR(64)
AS
IF @equipName IS NULL OR @userName IS NULL
    THROW 50001, 'Some Params is NULL', 1
DECLARE @userType VARCHAR(64)
EXEC usp_getUserUserType
@UName = @userName,
@UTName = @userType OUTPUT

IF @userType IS NULL OR @userType != 'Admin'
    THROW 50002, 'User Not Exist or the User is not admin', 1
BEGIN TRAN
INSERT INTO tblEquipment(equipName)
VALUES(@equipName)
IF @@ERROR<>0
    ROLLBACK TRAN
ELSE
    COMMIT TRAN

GO


CREATE PROCEDURE usp_getEquipID
@EName VARCHAR(128),
@EID INT OUTPUT
AS
SET @EID = (SELECT TOP 1 equipID FROM tblEquipment WHERE equipName = @EName)

GO



CREATE PROCEDURE usp_addEquipmentToRoom
@equipName VARCHAR(128),
@roomName VARCHAR(64),
@addDate DATE,
@userName VARCHAR(64)
AS
IF @equipName IS NULL OR @roomName IS NULL OR @addDate IS NULL OR @userName IS NULL
    THROW 50001, 'Some Params IS NULL', 1

DECLARE @roomID INT
EXEC usp_getRoomID
@RName = @roomName,
@RID = @roomID OUTPUT

IF @roomID IS NULL
    THROW 50002, 'Room NOT FOUND', 1

DECLARE @equipID INT
EXEC usp_getEquipID
@EName = @equipName,
@EID = @equipID OUTPUT

IF @equipID IS NULL
    THROW 50002, 'Equipment NOT FOUND', 1

DECLARE @userType VARCHAR(64)
EXEC usp_getUserUserType
@UName = @userName,
@UTName = @userType OUTPUT

IF @userType IS NULL OR @userType != 'Admin'
    THROW 50002, 'User is not admin', 1
BEGIN TRAN
INSERT INTO tblEquipInRoom(equipID, roomID, addDate)
VALUES(@equipID, @roomID, @addDate)

IF @@ERROR<>0
    ROLLBACK TRAN
ELSE
    COMMIT TRAN

GO

CREATE PROCEDURE usp_releaseReservation
@ResID INT,
@UserName VARCHAR(64)
AS
IF @ResID IS NULL OR @UserName IS NULL 
    THROW 50001, 'Some Params is null', 1

IF NOT EXISTS(SELECT * FROM tblReservation WHERE reservationID = @ResID)
    THROW 50002, 'Reservation not found', 1
DECLARE @targetUserID INT = (SELECT TOP 1 userID FROM tblReservation WHERE reservationID = @ResID)
DECLARE @givenUserID INT
EXEC usp_getUserID
@UName = @UserName,
@UID = @givenUserID OUTPUT
IF @givenUserID IS NULL OR @givenUserID != @targetUserID
    THROW 50002, 'Cur use either not exist or not the reservation owner', 1

BEGIN TRAN
DELETE FROM tblReservation WHERE reservationID = @ResID
IF @@ERROR<>0
    ROLLBACK TRAN
ELSE 
    COMMIT TRAN

GO

CREATE PROCEDURE usp_deleteRoom
@roomName VARCHAR(64),
@userName VARCHAR(64)
AS
IF @roomName IS NULL OR @userName IS NULL
    THROW 50001, 'Some PARAM is null', 1

DECLARE @userType VARCHAR(64)
EXEC usp_getUserUserType
@UName = @userName,
@UTName = @userType OUTPUT

IF @userType IS NULL OR @userType != 'Admin'
    THROW 50002, 'User is either not exist or not admin', 1

DECLARE @roomID INT
EXEC usp_getRoomID
@RName = @roomName,
@RID = @roomID OUTPUT

IF @roomID IS NULL
    THROW 50002, 'Room NOT FOUND', 1
BEGIN TRAN
DELETE FROM tblRoomIssue WHERE roomID = @roomID
DELETE FROM tblEquipInRoom WHERE roomID = @roomID
DELETE FROM tblReservation WHERE roomID = @roomID
DELETE FROM tblRoomStatusLog WHERE roomID = @roomID
DELETE FROM tblRoom WHERE roomID = @roomID

IF @@ERROR<>0
    ROLLBACK TRAN
ELSE
    COMMIT TRAN
GO

CREATE PROCEDURE usp_addIssue
@roomName VARCHAR(64),
@roomIssue VARCHAR(256),
@issueDate DATE
AS
IF @roomName IS NULL OR @roomIssue IS NULL OR @issueDate IS NULL
    THROW 50001, 'Some PARAM is null', 1
DECLARE @roomID INT
EXEC usp_getRoomID
@RName = @roomName,
@RID = @roomID OUTPUT

IF @roomID IS NULL
    THROW 50002, 'Room NOT FOUND', 1

BEGIN TRAN
INSERT INTO tblRoomIssue(roomID, roomIssueBody, createDate)
VALUES(@roomID, @roomIssue, @issueDate)
IF @@ERROR<>0
    ROLLBACK TRAN
ELSE
    COMMIT TRAN

GO

CREATE PROCEDURE usp_confirmIssue
@issueID INT,
@confirmDate DATE,
@userName VARCHAR(64)
AS
IF @issueID IS NULL OR @confirmDate IS NULL OR @userName IS NULL
    THROW 50001, 'Some PARAM is null', 1
IF NOT EXISTS(SELECT * FROM tblRoomIssue WHERE roomIssueID = @issueID)
    THROW 50002, 'Issue NOT FOUND', 1
DECLARE @hasConform DATE = (SELECT TOP 1 confirmDate FROM tblRoomIssue WHERE roomIssueID = @issueID)
IF @hasConform IS NOT NULL
    THROW 50002, 'Issue Has Conformed', 1
DECLARE @userType VARCHAR(64)
EXEC usp_getUserUserType
@UName = @userName,
@UTName = @userType OUTPUT

IF @userType IS NULL OR @userType != 'Admin'
    THROW 50002, 'User is either not exist or not admin', 1

BEGIN TRAN
UPDATE tblRoomIssue
SET confirmDate = @confirmDate
WHERE roomIssueID = @issueID
IF @@ERROR<>0
    ROLLBACK TRAN
ELSE
    COMMIT TRAN

GO

CREATE PROCEDURE usp_solveIssue
@issueID INT,
@solveDate DATE,
@userName VARCHAR(64)
AS
IF @issueID IS NULL OR @solveDate IS NULL OR @userName IS NULL
    THROW 50001, 'Some PARAM is null', 1
IF NOT EXISTS(SELECT * FROM tblRoomIssue WHERE roomIssueID = @issueID)
    THROW 50002, 'Issue NOT FOUND', 1
DECLARE @hasSolved DATE = (SELECT TOP 1 solveDate FROM tblRoomIssue WHERE roomIssueID = @issueID)
IF @hasSolved IS NOT NULL
    THROW 50002, 'Issue Has Solved', 1

DECLARE @userType VARCHAR(64)
EXEC usp_getUserUserType
@UName = @userName,
@UTName = @userType OUTPUT

IF @userType IS NULL OR @userType != 'Admin'
    THROW 50002, 'User is either not exist or not admin', 1
BEGIN TRAN
UPDATE tblRoomIssue
SET solveDate = @solveDate
WHERE roomIssueID = @issueID
IF @@ERROR<>0
    ROLLBACK TRAN
ELSE
    COMMIT TRAN
GO


CREATE PROCEDURE usp_removeEquipmentInRoom
@roomEquipID INT,
@userName VARCHAR(64),
@removeDate DATE
AS
IF @roomEquipID IS NULL OR @userName IS NULL OR @removeDate IS NULL
    THROW 50001, 'Some PARAM IS NULL', 1
IF NOT EXISTS(SELECT * FROM tblEquipInRoom WHERE equipRoomID = @roomEquipID)
    THROW 50002, 'Target Equipment is NOT in the room', 1
DECLARE @userType VARCHAR(64)
EXEC usp_getUserUserType
@UName = @userName,
@UTName = @userType OUTPUT

IF @userType IS NULL OR @userType != 'Admin'
    THROW 50002, 'User is either not exist or not admin', 1
BEGIN TRAN
UPDATE tblEquipInRoom
SET removeDate = @removeDate
WHERE equipRoomID = @roomEquipID

IF @@ERROR<>0
    ROLLBACK TRAN
ELSE
    COMMIT TRAN

GO

CREATE PROCEDURE usp_deleteEquipment
@equipName VARCHAR(128),
@userName VARCHAR(64)
AS
IF @equipName IS NULL OR @userName IS NULL
    THROW 50001, 'Some PARAM IS NULL', 1

DECLARE @equipID INT
EXEC usp_getEquipID
@EName = @equipName,
@EID = @equipID OUTPUT

IF @equipID IS NULL
    THROW 50002, 'Equipment NOT FOUND', 1

DECLARE @userType VARCHAR(64)
EXEC usp_getUserUserType
@UName = @userName,
@UTName = @userType OUTPUT

IF @userType IS NULL OR @userType != 'Admin'
    THROW 50002, 'User is either not exist or not admin', 1
BEGIN TRAN
DELETE FROM tblEquipInRoom WHERE equipID = @equipID
DELETE FROM tblEquipment WHERE equipID = @equipID
IF @@ERROR<>0
    ROLLBACK TRAN
ELSE
    COMMIT TRAN


GO

CREATE PROCEDURE usp_updateEquipmentName
@oldName VARCHAR(128),
@newName VARCHAR(128),
@userName VARCHAR(64)
AS
IF @oldName IS NULL OR @newName IS NULL OR @userName IS NULL
    THROW 50001, 'Some PARAM IS NULL', 1
DECLARE @equipID INT
EXEC usp_getEquipID
@EName = @oldName,
@EID = @equipID OUTPUT

IF @equipID IS NULL
    THROW 50001, 'Equipment NOT FOUND', 1

DECLARE @userType VARCHAR(64)
EXEC usp_getUserUserType
@UName = @userName,
@UTName = @userType OUTPUT

IF @userType IS NULL OR @userType != 'Admin'
    THROW 50002, 'User is either not exist or not admin', 1

BEGIN TRAN
UPDATE tblEquipment
SET equipName = @newName
WHERE equipID = @equipID
IF @@ERROR<>0
    ROLLBACK TRAN
ELSE
    COMMIT TRAN

GO

CREATE PROCEDURE usp_addNewUser
@userName VARCHAR(64),
@email VARCHAR(64),
@passHash BINARY(60),
@userTypeName VARCHAR(32)
AS
IF @userName IS NULL OR @email IS NULL OR @passHash IS NULL OR @userTypeName IS NULL
    THROW 50001, 'Some PARAM is null', 1
DECLARE @userTypeID INT
EXEC usp_getUserTypeID
@UTypename = @userTypeName,
@UTypeID = @userTypeID OUTPUT

IF @userTypeID IS NULL
    THROW 50002, 'User Type Not Found', 1
BEGIN TRAN
INSERT INTO tblUser(userName, email, passHash, userTypeID)
VALUES(@userName, @email, @passHash, @userTypeID)
IF @@ERROR<>0
    ROLLBACK TRAN
ELSE
    COMMIT TRAN

GO
