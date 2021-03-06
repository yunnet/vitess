<?php
// DO NOT EDIT! Generated by Protobuf-PHP protoc plugin 1.0
// Source: tabletmanagerdata.proto

namespace Vitess\Proto\Tabletmanagerdata {

  class GetPermissionsResponse extends \DrSlump\Protobuf\Message {

    /**  @var \Vitess\Proto\Tabletmanagerdata\Permissions */
    public $permissions = null;
    

    /** @var \Closure[] */
    protected static $__extensions = array();

    public static function descriptor()
    {
      $descriptor = new \DrSlump\Protobuf\Descriptor(__CLASS__, 'tabletmanagerdata.GetPermissionsResponse');

      // OPTIONAL MESSAGE permissions = 1
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 1;
      $f->name      = "permissions";
      $f->type      = \DrSlump\Protobuf::TYPE_MESSAGE;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $f->reference = '\Vitess\Proto\Tabletmanagerdata\Permissions';
      $descriptor->addField($f);

      foreach (self::$__extensions as $cb) {
        $descriptor->addField($cb(), true);
      }

      return $descriptor;
    }

    /**
     * Check if <permissions> has a value
     *
     * @return boolean
     */
    public function hasPermissions(){
      return $this->_has(1);
    }
    
    /**
     * Clear <permissions> value
     *
     * @return \Vitess\Proto\Tabletmanagerdata\GetPermissionsResponse
     */
    public function clearPermissions(){
      return $this->_clear(1);
    }
    
    /**
     * Get <permissions> value
     *
     * @return \Vitess\Proto\Tabletmanagerdata\Permissions
     */
    public function getPermissions(){
      return $this->_get(1);
    }
    
    /**
     * Set <permissions> value
     *
     * @param \Vitess\Proto\Tabletmanagerdata\Permissions $value
     * @return \Vitess\Proto\Tabletmanagerdata\GetPermissionsResponse
     */
    public function setPermissions(\Vitess\Proto\Tabletmanagerdata\Permissions $value){
      return $this->_set(1, $value);
    }
  }
}

