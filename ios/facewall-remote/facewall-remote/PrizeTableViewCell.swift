//
//  PrizeTableViewCell.swift
//  facewall-remote
//
//  Created by Lingfei Song on 11/7/15.
//  Copyright © 2015 zealion. All rights reserved.
//

import UIKit

class PrizeTableViewCell: UITableViewCell {

    @IBOutlet weak var prizeLabel: UILabel!
    @IBOutlet weak var sumLabel: UILabel!
    override func awakeFromNib() {
        super.awakeFromNib()
        // Initialization code
    }

    override func setSelected(selected: Bool, animated: Bool) {
        super.setSelected(selected, animated: animated)

        // Configure the view for the selected state
    }

}
