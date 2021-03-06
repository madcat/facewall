//
//  AppDelegate.swift
//  facewall-remote
//
//  Created by Lingfei Song on 11/7/15.
//  Copyright © 2015 zealion. All rights reserved.
//

import UIKit

@UIApplicationMain
class AppDelegate: UIResponder, UIApplicationDelegate {

    var window: UIWindow?


    func application(application: UIApplication, didFinishLaunchingWithOptions launchOptions: [NSObject: AnyObject]?) -> Bool {
        registerDefaultsFromSettingsBundle()
        return true
    }

    func applicationWillResignActive(application: UIApplication) {
        // Sent when the application is about to move from active to inactive state. This can occur for certain types of temporary interruptions (such as an incoming phone call or SMS message) or when the user quits the application and it begins the transition to the background state.
        // Use this method to pause ongoing tasks, disable timers, and throttle down OpenGL ES frame rates. Games should use this method to pause the game.
    }

    func applicationDidEnterBackground(application: UIApplication) {
        // Use this method to release shared resources, save user data, invalidate timers, and store enough application state information to restore your application to its current state in case it is terminated later.
        // If your application supports background execution, this method is called instead of applicationWillTerminate: when the user quits.
    }

    func applicationWillEnterForeground(application: UIApplication) {
        // Called as part of the transition from the background to the inactive state; here you can undo many of the changes made on entering the background.
    }

    func applicationDidBecomeActive(application: UIApplication) {
        // Restart any tasks that were paused (or not yet started) while the application was inactive. If the application was previously in the background, optionally refresh the user interface.
    }

    func applicationWillTerminate(application: UIApplication) {
        // Called when the application is about to terminate. Save data if appropriate. See also applicationDidEnterBackground:.
    }

    func registerDefaultsFromSettingsBundle() {
        let defs = NSUserDefaults.standardUserDefaults()
        defs.synchronize()
        
        guard let path = NSBundle.mainBundle().pathForResource("Settings", ofType:"bundle") else {
            print("could not find settings.bundle")
            return
        }
        
        guard let dict = NSDictionary(contentsOfFile:NSString(string:path).stringByAppendingPathComponent("Root.plist")), prefs = dict.objectForKey("PreferenceSpecifiers") as? Array<Dictionary<String,AnyObject>> else {
            print("incorrect Settings.bundle format")
            return
        }
        
        var defaultsToRegister = [String:AnyObject]()
        for pref in prefs {
            if let key = pref["Key"] as? String where defs.objectForKey(key) == nil {
                defaultsToRegister[key] = pref["DefaultValue"]
            }
        }
        
        defs.registerDefaults(defaultsToRegister)
        defs.synchronize()
        
    }
}

