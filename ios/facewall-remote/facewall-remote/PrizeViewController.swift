//
//  ViewController.swift
//  facewall-remote
//
//  Created by Lingfei Song on 11/7/15.
//  Copyright © 2015 zealion. All rights reserved.
//

import UIKit

class PrizeViewController: UIViewController {

    override func viewDidLoad() {
        super.viewDidLoad()

        guard let host = NSUserDefaults.standardUserDefaults().stringForKey("host") else {
            var alert = UIAlertController(title: "错误", message: "未设置服务器地址，请设置并重启应用", preferredStyle: UIAlertControllerStyle.Alert)
            //alert.addAction(UIAlertAction(title: "Click", style: UIAlertActionStyle.Default, handler: nil))
            self.presentViewController(alert, animated: true, completion: nil)
            return
        }
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }


}

