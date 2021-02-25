import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RegistrySettingComponent } from './registry-setting.component';

describe('RegistrySettingComponent', () => {
  let component: RegistrySettingComponent;
  let fixture: ComponentFixture<RegistrySettingComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RegistrySettingComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(RegistrySettingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
